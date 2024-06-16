package service

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"nexa/services/file_storage/internal/domain/dto"
	"nexa/services/file_storage/internal/domain/external"
	"nexa/services/file_storage/internal/domain/mapper"
	"nexa/services/file_storage/internal/domain/repository"
	"nexa/services/file_storage/internal/domain/service"
	"nexa/services/file_storage/util"
	spanUtil "nexa/shared/span"
	"nexa/shared/status"
	"nexa/shared/types"
)

func NewFileStorage(metadataRepo repository.IFileMetadata, storage external.IStorage) service.IFileStorage {
	return &fileStorage{
		metadataRepo: metadataRepo,
		storageExt:   storage,
		tracer:       util.GetTracer(),
	}
}

type fileStorage struct {
	metadataRepo repository.IFileMetadata
	storageExt   external.IStorage

	tracer trace.Tracer
}

func (f fileStorage) Store(ctx context.Context, file *dto.FileStoreDTO) (string, status.Object) {
	ctx, span := f.tracer.Start(ctx, "FileStorageService.Store")
	defer span.End()

	// upload to storage
	files := mapper.MapFileStoreDTO(file)
	relativePath, err := f.storageExt.Store(ctx, &files)
	if err != nil {
		spanUtil.RecordError(err, span)
		return "", status.New(status.EXTERNAL_SERVICE_ERROR, err)
	}

	// save metadata
	metadata, err := mapper.MapFileMetadata(file, f.storageExt)
	metadata.ProviderPath = relativePath
	if err != nil {
		spanUtil.RecordError(err, span)
		return "", status.ErrInternal(err)
	}

	err = f.metadataRepo.Create(ctx, &metadata)
	if err != nil {
		spanUtil.RecordError(err, span)
		return "", status.FromRepository(err, status.NullCode)
	}

	return metadata.Id.Underlying().String(), status.Created()
}

func (f fileStorage) Find(ctx context.Context, id types.Id) (dto.FileResponseDTO, status.Object) {
	ctx, span := f.tracer.Start(ctx, "FileStorageService.Find")
	defer span.End()

	// get metadata
	metadata, err := f.metadataRepo.FindByIds(ctx, id)
	if err != nil {
		spanUtil.RecordError(err, span)
		return dto.FileResponseDTO{}, status.FromRepository(err, status.NullCode)
	}

	// get the file
	file, err := f.storageExt.Find(ctx, metadata[0].ProviderPath)
	if err != nil {
		spanUtil.RecordError(err, span)
		return dto.FileResponseDTO{}, status.New(status.EXTERNAL_SERVICE_ERROR, err)
	}

	return mapper.ToFileResponse(&file), status.Success()
}

func (f fileStorage) FindMetadata(ctx context.Context, id types.Id) (*dto.FileMetadataResponseDTO, status.Object) {
	ctx, span := f.tracer.Start(ctx, "FileStorageService.FindMetadata")
	defer span.End()

	// get metadata
	metadata, err := f.metadataRepo.FindByIds(ctx, id)
	if err != nil {
		spanUtil.RecordError(err, span)
		return nil, status.FromRepository(err, status.NullCode)
	}

	resp := mapper.ToFileMetadataResponse(&metadata[0]) // heap allocated
	return &resp, status.Success()
}

func (f fileStorage) Delete(ctx context.Context, id types.Id) status.Object {
	ctx, span := f.tracer.Start(ctx, "FileStorageService.Delete")
	defer span.End()

	// get metadata
	metadata, err := f.metadataRepo.FindByIds(ctx, id)
	if err != nil {
		spanUtil.RecordError(err, span)
		return status.FromRepository(err, status.NullCode)
	}

	// delete from storage
	err = f.storageExt.Delete(ctx, metadata[0].ProviderPath)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}

	// delete metadata
	err = f.metadataRepo.DeleteById(ctx, id)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}

	return status.Deleted()
}

func (f fileStorage) UpdateMetadata(ctx context.Context, input *dto.UpdateFileMetadataDTO) status.Object {
	ctx, span := f.tracer.Start(ctx, "FileStorageService.UpdateMetadata")
	defer span.End()

	obj, err := mapper.MapUpdateFileMetadataDTO(input)
	if err != nil {
		return status.ErrInternal(err)
	}

	err = f.metadataRepo.Patch(ctx, &obj)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}

	return status.Updated()
}
