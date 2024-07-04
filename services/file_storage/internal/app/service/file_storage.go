package service

import (
  "context"
  "fmt"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/file_storage/internal/app/uow"
  "nexa/services/file_storage/internal/domain/dto"
  "nexa/services/file_storage/internal/domain/external"
  "nexa/services/file_storage/internal/domain/mapper"
  "nexa/services/file_storage/internal/domain/service"
  "nexa/services/file_storage/util"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUow "nexa/shared/uow"
  spanUtil "nexa/shared/util/span"
)

func NewFileStorage(unit sharedUow.IUnitOfWork[uow.FileMetadataStorage], storage external.IStorage) service.IFileStorage {
  return &fileStorage{
    unit:       unit,
    storageExt: storage,
    tracer:     util.GetTracer(),
  }
}

type fileStorage struct {
  unit       sharedUow.IUnitOfWork[uow.FileMetadataStorage]
  storageExt external.IStorage

  tracer trace.Tracer
}

func (f *fileStorage) publicPath(filename string) string {
  return fmt.Sprintf("/public/%s", filename)
}

func (f *fileStorage) Store(ctx context.Context, fileDto *dto.FileStoreDTO) (types.Id, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.Store")
  defer span.End()

  // Map to domain
  file, metadata, err := fileDto.ToDomain(f.storageExt)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  // Upload file
  relativePath, err := f.storageExt.Store(ctx, &file)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrExternal(err)
  }

  // Save metadata
  metadata.ProviderPath = relativePath
  repos := f.unit.Repositories()
  err = repos.Metadata().Create(ctx, &metadata)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  return metadata.Id, status.Created()
}

func (f *fileStorage) Find(ctx context.Context, id types.Id) (dto.FileResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.Find")
  defer span.End()

  // get metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.FileResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  // get the file
  file, err := f.storageExt.Find(ctx, metadata[0].ProviderPath)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.FileResponseDTO{}, status.ErrExternal(err)
  }

  return mapper.ToFileResponse(&file), status.Success()
}

func (f *fileStorage) FindMetadata(ctx context.Context, id types.Id) (*dto.FileMetadataResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.FindMetadata")
  defer span.End()

  // get metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  resp := mapper.ToFileMetadataResponse(&metadata[0]) // heap allocated
  return &resp, status.Success()
}

func (f *fileStorage) Delete(ctx context.Context, id types.Id) status.Object {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.Delete")
  defer span.End()

  // Get metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  stat := status.Deleted()
  // err ignored, because it could be caused by repository and storage external service
  _ = f.unit.DoTx(ctx, func(ctx context.Context, storage uow.FileMetadataStorage) error {
    ctx, span := f.tracer.Start(ctx, "UOW.Delete")
    // Delete metadata from persistent database
    err = storage.Metadata().DeleteById(ctx, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete from storage
    err = f.storageExt.Delete(ctx, metadata[0].ProviderPath)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    return nil
  })

  return stat
}

func (f *fileStorage) UpdateMetadata(ctx context.Context, input *dto.UpdateFileMetadataDTO) status.Object {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.UpdateMetadata")
  defer span.End()

  obj := input.ToDomain()
  // Get file metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, input.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Doesn't need to update, currently it only handle moving into/from public path
  if metadata[0].IsPublic == obj.IsPublic {
    return status.Updated()
  }

  stat := status.Updated()
  _ = f.unit.DoTx(ctx, func(ctx context.Context, storage uow.FileMetadataStorage) error {
    ctx, span = f.tracer.Start(ctx, "UOW.UpdateMetadata")
    defer span.End()

    // Copy file
    dest := metadata[0].Name
    if obj.IsPublic {
      dest = f.publicPath(dest)
    }
    newPath, err := f.storageExt.Copy(ctx, metadata[0].ProviderPath, dest)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    // Update file metadata
    obj.ProviderPath = newPath
    err = storage.Metadata().Update(ctx, &obj)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete old file
    err = f.storageExt.Delete(ctx, metadata[0].ProviderPath)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    return nil
  })

  return stat
}
