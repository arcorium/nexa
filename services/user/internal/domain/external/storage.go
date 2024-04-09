package external

import (
  "context"
  "nexa/shared/types"
)

type IStorageClient interface {
  UploadImage(ctx context.Context) (types.FilePath, error)
}
