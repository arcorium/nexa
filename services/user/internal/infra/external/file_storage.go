package external

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  storagev1 "nexa/proto/gen/go/file_storage/v1"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/external"
  "nexa/services/user/util"
  "nexa/shared/types"
  spanUtil "nexa/shared/util/span"
)

func NewFileStorageClient(conn grpc.ClientConnInterface) external.IFileStorageClient {
  return &fileStorageClient{
    client: storagev1.NewFileStorageServiceClient(conn),
    tracer: util.GetTracer(),
  }
}

type fileStorageClient struct {
  client storagev1.FileStorageServiceClient
  tracer trace.Tracer
}

func (f *fileStorageClient) UploadProfileImage(ctx context.Context, dto *dto.UploadImageDTO) (types.Id, types.FilePath, error) {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.UploadProfileImage")
  defer span.End()

  stream, err := f.client.Store(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), "", err
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
      spanUtil.RecordError(err, span)
      return types.NullId(), "", err
    }
    break
  }

  recv, err := stream.CloseAndRecv()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), "", err
  }

  id, err := types.IdFromString(recv.FileId)

  return id, types.FilePathFromString(*recv.Filepath), err
}

func (f *fileStorageClient) DeleteProfileImage(ctx context.Context, id types.Id) error {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.DeleteProfileImage")
  defer span.End()

  req := storagev1.DeleteFileRequest{FileId: id.String()}
  _, err := f.client.Delete(ctx, &req)
  return err
}
