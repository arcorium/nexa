package handler

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/grpc/interceptor"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "io"
  "nexa/services/file_storage/internal/api/grpc/mapper"
  "nexa/services/file_storage/internal/domain/dto"
  "nexa/services/file_storage/internal/domain/service"
  "nexa/services/file_storage/util"
)

func NewFileStorage(file service.IFileStorage) StorageHandler {
  return StorageHandler{
    fileService: file,
    tracer:      util.GetTracer(),
  }
}

type StorageHandler struct {
  storagev1.FileStorageServiceServer
  fileService service.IFileStorage

  tracer trace.Tracer
}

func (s *StorageHandler) Register(server *grpc.Server) {
  storagev1.RegisterFileStorageServiceServer(server, s)
}

func (s *StorageHandler) Find(request *storagev1.FindFileRequest, server storagev1.FileStorageService_FindServer) error {
  stream, err := interceptor.GetWrappedServerStream(server)
  ctx := server.Context()
  if err == nil {
    ctx = stream.WrappedContext
  }

  ctx, span := s.tracer.Start(ctx, "StorageHandler.Find")
  defer span.End()

  id, err := types.IdFromString(request.FileId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedErr.NewFieldError("file_id", err).ToGrpcError()
  }

  file, stat := s.fileService.Find(ctx, id)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat.ToGRPCError()
  }

  // TODO: Split bytes into several chunks for large file
  for {
    err = server.Send(&storagev1.FindFileResponse{
      Filename: file.Name,
      Chunk:    file.Data,
    })

    if err != nil {
      spanUtil.RecordError(err, span)
      return err
    }
    break
  }

  return nil
}

func (s *StorageHandler) FindMetadata(ctx context.Context, request *storagev1.FindFileMetadataRequest) (*storagev1.FindFileMetadataResponse, error) {
  ctx, span := s.tracer.Start(ctx, "StorageHandler.FindMetadata")
  defer span.End()

  fileId, err := types.IdFromString(request.FileId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("file_id", err).ToGrpcError()
  }

  metadata, stat := s.fileService.FindMetadata(ctx, fileId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := &storagev1.FindFileMetadataResponse{File: mapper.ToProtoFile(metadata)}
  return response, nil
}

func (s *StorageHandler) Store(server storagev1.FileStorageService_StoreServer) error {
  stream, err := interceptor.GetWrappedServerStream(server)
  ctx := server.Context()
  if err == nil {
    ctx = stream.WrappedContext
  }

  ctx, span := s.tracer.Start(ctx, "StorageHandler.FindMetadata")
  defer span.End()

  storeDto := dto.FileStoreDTO{}
  for {
    req, err := server.Recv()
    if err != nil {
      if err == io.EOF {
        break
      }
      spanUtil.RecordError(err, span)
      return err
    }
    storeDto.Name = req.Filename
    storeDto.IsPublic = req.IsPublic
    storeDto.Data = append(storeDto.Data, req.Chunk...)
  }

  // Validate
  err = sharedUtil.ValidateStruct(&storeDto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  id, stat := s.fileService.Store(ctx, &storeDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    server.SendAndClose(nil)
    return stat.ToGRPCError()
  }
  return server.SendAndClose(&storagev1.StoreFileResponse{FileId: id.String()})
}

func (s *StorageHandler) Update(ctx context.Context, request *storagev1.UpdateFileRequest) (*emptypb.Empty, error) {
  ctx, span := s.tracer.Start(ctx, "StorageHandler.Update")
  defer span.End()

  dtos, err := mapper.ToUpdateMetadataDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := s.fileService.Move(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (s *StorageHandler) Delete(ctx context.Context, request *storagev1.DeleteFileRequest) (*emptypb.Empty, error) {
  ctx, span := s.tracer.Start(ctx, "StorageHandler.Delete")
  defer span.End()

  id, err := types.IdFromString(request.FileId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("file_id", err).ToGrpcError()
  }

  stat := s.fileService.Delete(ctx, id)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
