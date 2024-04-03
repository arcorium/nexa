package external

import (
	"context"
	"nexa/shared/types"
)

type IFileStorageClient interface {
	UploadImage(ctx context.Context) (types.FilePath, error)
}
