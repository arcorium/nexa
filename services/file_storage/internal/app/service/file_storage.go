package service

import (
  "context"
  "errors"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUow "github.com/arcorium/nexa/shared/uow"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/file_storage/constant"
  "nexa/services/file_storage/internal/app/uow"
  "nexa/services/file_storage/internal/domain/dto"
  "nexa/services/file_storage/internal/domain/external"
  "nexa/services/file_storage/internal/domain/mapper"
  "nexa/services/file_storage/internal/domain/service"
  "nexa/services/file_storage/util"
)

func NewFileStorage(unit sharedUow.IUnitOfWork[uow.FileMetadataStorage], storage external.IStorage) service.IFileStorage {
  return &fileStorageService{
    unit:       unit,
    storageExt: storage,
    tracer:     util.GetTracer(),
  }
}

type fileStorageService struct {
  unit       sharedUow.IUnitOfWork[uow.FileMetadataStorage]
  storageExt external.IStorage

  tracer trace.Tracer
}

func (f *fileStorageService) Store(ctx context.Context, storeDTO *dto.FileStoreDTO) (dto.FileStoreResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.Store")
  defer span.End()

  // Map to domain
  file, metadata, err := storeDTO.ToDomain(f.storageExt.GetProvider())
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.FileStoreResponseDTO{}, status.ErrInternal(err)
  }

  // Check permissions
  if !storeDTO.IsPublic {
    claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
    if !authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_STORE_PRIVATE]) {
      return dto.FileStoreResponseDTO{}, status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
    }
  }

  var stat = status.Success()
  err = f.unit.DoTx(ctx, func(ctx context.Context, storage uow.FileMetadataStorage) error {
    ctx, span := f.tracer.Start(ctx, "UOW.Store")
    defer span.End()
    // Upload file
    relativePath, err := f.storageExt.Store(ctx, &file)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    defer func() {
      // Delete file when there is an error on get full path or creating metadata
      if err == nil {
        return
      }
      // Delete file from storage
      err2 := f.storageExt.Delete(ctx, relativePath)
      if err2 != nil {
        spanUtil.RecordError(err2, span)
        stat = status.ErrExternal(err2)
      }
    }()

    // Save metadata
    metadata.ProviderPath = relativePath

    // Get fullpath for public file
    if storeDTO.IsPublic {
      path, err := f.storageExt.GetFullPath(ctx, relativePath)
      if err != nil {
        spanUtil.RecordError(err, span)
        stat = status.ErrExternal(err)
        return err
      }
      metadata.FullPath = path.Path()
    }

    err = storage.Metadata().Create(ctx, &metadata)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    return nil
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.FileStoreResponseDTO{}, stat
  }

  resp := dto.FileStoreResponseDTO{
    Id:       metadata.Id,
    FullPath: types.FilePathFromString(metadata.FullPath),
  }
  return resp, stat
}

func (f *fileStorageService) Find(ctx context.Context, id types.Id) (dto.FileResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.Find")
  defer span.End()

  // get metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.FileResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  // Check permission to get private file
  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  if !metadata[0].IsPublic && !authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_GET]) {
    err = sharedErr.ErrUnauthorizedPermission
    spanUtil.RecordError(err, span)
    return dto.FileResponseDTO{}, status.ErrUnAuthorized(err)
  }

  // get the file
  file, err := f.storageExt.Find(ctx, metadata[0].ProviderPath)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.FileResponseDTO{}, status.ErrExternal(err)
  }

  // set original name as name
  file.Name = metadata[0].Name
  return mapper.ToFileResponse(&file), status.Success()
}

func (f *fileStorageService) FindMetadatas(ctx context.Context, ids ...types.Id) ([]dto.FileMetadataResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.FindMetadatas")
  defer span.End()

  // get metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, ids...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(metadata, mapper.ToFileMetadataResponse) // heap allocated
  return resp, status.Success()
}

func (f *fileStorageService) Delete(ctx context.Context, id types.Id) status.Object {
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
    defer span.End()
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

func (f *fileStorageService) Move(ctx context.Context, updateDto *dto.UpdateFileMetadataDTO) status.Object {
  ctx, span := f.tracer.Start(ctx, "FileStorageService.Move")
  defer span.End()

  // Get file metadata
  repos := f.unit.Repositories()
  metadata, err := repos.Metadata().FindByIds(ctx, updateDto.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Doesn't need to update, currently it only handle moving into/from public path
  if metadata[0].IsPublic == updateDto.IsPublic {
    return status.New(status.OBJECT_NOT_FOUND, errors.New("nothings to do, the file is already on right location"))
  }

  stat := status.Updated()
  _ = f.unit.DoTx(ctx, func(ctx context.Context, storage uow.FileMetadataStorage) error {
    ctx, span = f.tracer.Start(ctx, "UOW.Move")
    defer span.End()

    // Get destination path
    dest := f.storageExt.GetProviderPath(metadata[0].ProviderPath, updateDto.IsPublic)

    // Copy file
    newPath, err := f.storageExt.Copy(ctx, metadata[0].ProviderPath, dest.String())
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    defer func() {
      // When there in an error, delete copied file
      if err == nil {
        return
      }
      // Delete copied file
      err2 := f.storageExt.Delete(ctx, newPath)
      if err2 != nil {
        spanUtil.RecordError(err2, span)
        stat = status.ErrExternal(err2)
      }
    }()

    // Patch file metadata
    patched := updateDto.ToDomain(newPath)

    // Get fullpath to access directly
    if updateDto.IsPublic {
      fullpath, err := f.storageExt.GetFullPath(ctx, newPath)
      if err != nil {
        spanUtil.RecordError(err, span)
        stat = status.ErrExternal(err)
        return err
      }
      patched.FullPath = types.SomeNullable(fullpath.Path())
    }

    err = storage.Metadata().Patch(ctx, &patched)
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
