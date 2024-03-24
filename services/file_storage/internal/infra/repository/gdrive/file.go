package gdrive

import (
	"nexa/services/file_storage/internal/domain/repository"
	"nexa/services/file_storage/shared/domain/entity"
	"nexa/shared/types"
)

func NewFileRepository() repository.IFile {
	return &fileRepository{}
}

type fileRepository struct {
}

func (f fileRepository) Store(filename string, bytes []byte) (types.FilePath, error) {
	//TODO implement me
	panic("implement me")
}

func (f fileRepository) Find(filename string) (entity.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f fileRepository) Delete(filename string, bin bool) error {
	//TODO implement me
	panic("implement me")
}

func (f fileRepository) Restore(filename string) error {
	//TODO implement me
	panic("implement me")
}
