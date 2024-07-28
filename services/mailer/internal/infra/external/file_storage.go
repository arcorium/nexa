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
  "io"
  "nexa/services/mailer/config"
  "nexa/services/mailer/internal/domain/dto"
  "nexa/services/mailer/internal/domain/external"
  "nexa/services/mailer/util"
)

func NewFileStorageClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IFileStorageClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-media-storage",
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

func (f *fileStorageClient) GetFiles(ctx context.Context, fileIds ...types.Id) ([]dto.FileAttachment, error) {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.GetFiles")
  defer span.End()

  var result []dto.FileAttachment
  _, err := f.cb.Execute(func() (interface{}, error) {
    for _, id := range fileIds {
      // TODO: Allow empty file
      // Init connection
      find, err := f.client.Find(ctx, &storagev1.FindFileRequest{FileId: id.String()})
      if err != nil {
        return nil, err
      }
      // Receive stream
      current := dto.FileAttachment{}
      for {
        recv, err := find.Recv()
        if err != nil {
          if err == io.EOF {
            result = append(result, current)
            break
          }
          // Skip file that doesn't exist
          s, ok := status.FromError(err)
          if !ok || s.Code() != codes.NotFound {
            return nil, err
          }
          break
        }
        current.Filename = recv.Filename
        current.Data = append(current.Data, recv.Chunk...)
      }
    }
    return nil, nil
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return nil, err
  }

  return result, nil
}
