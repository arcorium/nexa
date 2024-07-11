package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  domain "nexa/services/file_storage/internal/domain/entity"
)

type IFileMetadata interface {
  FindByIds(ctx context.Context, ids ...types.Id) ([]domain.FileMetadata, error)
  FindByNames(ctx context.Context, names ...string) ([]domain.FileMetadata, error)
  Create(ctx context.Context, metadata *domain.FileMetadata) error
  Patch(ctx context.Context, metadata *domain.PatchedFileMetadata) error
  DeleteById(ctx context.Context, id types.Id) error
  DeleteByName(ctx context.Context, name string) error
}
