package external

import "nexa/shared/types"

type IFileStorageClient interface {
	UploadImage() (types.FilePath, error)
}
