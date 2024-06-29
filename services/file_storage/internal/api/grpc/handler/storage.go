package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "io"
  protoV1 "nexa/proto/generated/golang/file_storage/v1"
  "nexa/services/file_storage/internal/api/grpc/mapper"
  "nexa/services/file_storage/internal/domain/dto"
  "nexa/services/file_storage/internal/domain/service"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func NewFileStorage(file service.IFileStorage) StorageHandler {
  return StorageHandler{
    fileService: file,
  }
}

type StorageHandler struct {
  protoV1.UnimplementedFileStorageServiceServer

  fileService service.IFileStorage
}

func (s *StorageHandler) Register(server *grpc.Server) {
  protoV1.RegisterFileStorageServiceServer(server, s)
}

func (s *StorageHandler) Find(request *protoV1.FindFileRequest, server protoV1.FileStorageService_FindServer) error {
  span := trace.SpanFromContext(server.Context())

  id, err := types.IdFromString(request.FileId)
  if err != nil {
    spanUtil.RecordError(err, span)
    stat := status.ErrBadRequest(err)
    return stat.ToGRPCError()
  }

  file, stat := s.fileService.Find(server.Context(), id)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat.ToGRPCError()
  }

  // TODO: Split bytes into several chunks for large file
  for {
    err = server.Send(&protoV1.FindFileResponse{
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

func (s *StorageHandler) FindMetadata(ctx context.Context, request *protoV1.FindFileMetadataRequest) (*protoV1.FindFileMetadataResponse, error) {
  span := trace.SpanFromContext(ctx)

  id, err := types.IdFromString(request.FileId)
  if err != nil {
    spanUtil.RecordError(err, span)
    stat := status.ErrBadRequest(err)
    return nil, stat.ToGRPCError()
  }

  metadata, stat := s.fileService.FindMetadata(ctx, id)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := &protoV1.FindFileMetadataResponse{File: mapper.ToProtoFile(metadata)}
  return response, nil
}

func (s *StorageHandler) Store(server protoV1.FileStorageService_StoreServer) error {
  span := trace.SpanFromContext(server.Context())

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
    storeDto.Data = append(storeDto.Data, req.Chunk...)
  }
  id, stat := s.fileService.Store(server.Context(), &storeDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    server.SendAndClose(nil)
    return stat.ToGRPCError()
  }
  return server.SendAndClose(&protoV1.StoreFileResponse{FileId: id})
}

func (s *StorageHandler) UpdateMetadata(ctx context.Context, request *protoV1.UpdateFileMetadataRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  input := mapper.ToUpdateMetadataDTO(request)
  err := sharedUtil.ValidateStructCtx(ctx, &input)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := s.fileService.UpdateMetadata(ctx, &input)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return &emptypb.Empty{}, stat.ToGRPCError()
}

func (s *StorageHandler) Delete(ctx context.Context, request *protoV1.DeleteFileRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  id, err := types.IdFromString(request.FileId)
  if err != nil {
    spanUtil.RecordError(err, span)
    stat := status.ErrBadRequest(err)
    return nil, stat.ToGRPCError()
  }

  stat := s.fileService.Delete(ctx, id)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}
