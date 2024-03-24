package service

import (
	"context"
	"nexa/services/file_storage/shared/domain/entity"
	"nexa/shared/status"
	"nexa/shared/types"
)

type IFile interface {
	// Store store file
	Store(ctx context.Context) (types.FilePath, status.Object)
	// Load read file based on the filename
	Load(ctx context.Context) (entity.File, status.Object)
	// Delete remove or place file on bin based on the filename
	Delete(ctx context.Context) status.Object
	// Replace delete old file and upload new file with the same name
	Replace(ctx context.Context) (types.FilePath, status.Object)
	// Restore restore file on bin
	Restore(ctx context.Context) status.Object
}
