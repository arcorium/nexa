package external

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  storagev1 "nexa/proto/gen/go/file_storage/v1"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/external"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
)

func NewFileStorageClient(conn grpc.ClientConnInterface, tracer trace.Tracer) external.IFileStorageClient {
  return &fileStorageClient{
    client: storagev1.NewFileStorageServiceClient(conn),
    tracer: tracer,
  }
}

type fileStorageClient struct {
  client storagev1.FileStorageServiceClient
  tracer trace.Tracer
}

func (f *fileStorageClient) upload(ctx context.Context, filename string, data []byte) (types.Id, error) {
  span := trace.SpanFromContext(ctx)

  stream, err := f.client.Store(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  for {
    request := storagev1.StoreFileRequest{
      Filename: filename,
      IsPublic: true,
      Chunk:    data,
    }

    err = stream.Send(&request)
    if err != nil {
      spanUtil.RecordError(err, span)
      return types.NullId(), err
    }
    break
  }

  recv, err := stream.CloseAndRecv()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  id, err := types.IdFromString(recv.FileId)
  return id, err
}

func (f *fileStorageClient) UploadProfileImage(ctx context.Context, dto *dto.UploadImageDTO) (types.Id, error) {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.UploadProfileImage")
  defer span.End()

  return f.upload(ctx, dto.Filename, dto.Data)
}

func (f *fileStorageClient) UpdateProfileImage(ctx context.Context, dto *dto.UpdateImageDTO) (types.Id, error) {
  ctx, span := f.tracer.Start(ctx, "FileStorageClient.UpdateProfileImage")
  defer span.End()

  deleteRequest := storagev1.DeleteFileRequest{FileId: dto.Id.Underlying().String()}

  // Delete
  _, err := f.client.Delete(ctx, &deleteRequest)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  // Upload new one
  return f.upload(ctx, dto.Filename, dto.Data)
}
