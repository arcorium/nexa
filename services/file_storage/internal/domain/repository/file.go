package repository

import (
	"nexa/services/file_storage/shared/domain/entity"
	"nexa/shared/types"
)

type IFile interface {
	Store(filename string, bytes []byte) (types.FilePath, error)
	Find(filename string) (entity.File, error)
	Delete(filename string, bin bool) error
	Restore(filename string) error
}
