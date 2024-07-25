package external

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "nexa/services/authentication/config"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
)

func NewFileStorageClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IFileStorageClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-file_storage",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &fileStorageClient{
    client: storagev1.NewFileStorageServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type fileStorageClient struct {
  client storagev1.FileStorageServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (f *fileStorageClient) UploadProfileImage(ctx context.Context, dto *dto.UploadImageDTO) (types.Id, types.FilePath, error) {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.UploadProfileImage")
  defer span.End()

  result, err := f.cb.Execute(func() (interface{}, error) {
    stream, err := f.client.Store(ctx)
    if err != nil {
      return nil, err
    }

    for {
      request := storagev1.StoreFileRequest{
        Filename: dto.Filename,
        IsPublic: true,
        Chunk:    dto.Data,
      }

      // TODO: Split data
      err = stream.Send(&request)
      if err != nil {
        return nil, err
      }
      break
    }

    return stream.CloseAndRecv()
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), "", err
  }

  resp := result.(*storagev1.StoreFileResponse)
  id, err := types.IdFromString(resp.FileId)
  path := types.NewNullable(resp.Filepath)
  return id, types.FilePathFromString(path.ValueOr("")), nil
}

func (f *fileStorageClient) DeleteProfileImage(ctx context.Context, id types.Id) error {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.DeleteProfileImage")
  defer span.End()

  req := storagev1.DeleteFileRequest{FileId: id.String()}
  _, err := f.cb.Execute(func() (interface{}, error) {
    return f.client.Delete(ctx, &req)
  })
  stat, ok := status.FromError(err)
  if !ok {
    return err
  }
  // Allow not found status
  if stat.Code() == codes.NotFound {
    return nil
  }
  return err
}
