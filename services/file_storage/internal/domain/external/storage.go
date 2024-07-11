package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  domain "nexa/services/file_storage/internal/domain/entity"
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
  // GetProviderPath will return the path for expected location. it should not call any external client,
  // instead it should just create new path. for example for filename.png as filename and public is true
  // it should return something like public/filename.png or only filename.png for public is false
  GetProviderPath(filename string, public bool) types.FilePath
  GetProvider() domain.StorageProvider
}
