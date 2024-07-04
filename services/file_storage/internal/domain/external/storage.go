package external

import (
  "context"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/shared/types"
)

type IStorage interface {
  Close(ctx context.Context) error
  Find(ctx context.Context, filename string) (domain.File, error)
  // Store upload files on bucket and returning the id that could be the path or the object/file id.
  // The id/path can be used for another operation
  Store(ctx context.Context, file *domain.File) (string, error)
  // Copy file into another location
  Copy(ctx context.Context, src, dest string) (string, error)
  Delete(ctx context.Context, filename string) error
  // GetFullPath get download path for specific file
  GetFullPath(ctx context.Context, filename string) (types.FilePath, error)
  GetProvider() domain.StorageProvider
}
