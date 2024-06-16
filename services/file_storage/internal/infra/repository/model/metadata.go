package model

import (
	"github.com/uptrace/bun"
	domain "nexa/services/file_storage/internal/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
	"nexa/shared/variadic"
	"nexa/shared/wrapper"
	"time"
)

type FileMapOption = repo.DataAccessModelMapOption[*domain.FileMetadata, *FileMetadata]

func FromFileDomain(domain *domain.FileMetadata, opts ...FileMapOption) FileMetadata {
	obj := FileMetadata{
		Id:              domain.Id.Underlying().String(),
		Filename:        domain.Name,
		FileType:        domain.Type,
		Size:            domain.Size,
		IsPublic:        domain.IsPublic,
		StorageProvider: domain.Provider.Underlying(),
		StoragePath:     domain.ProviderPath,
	}

	variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &obj))
	return obj
}

type FileMetadata struct {
	bun.BaseModel `bun:"table:file_metadata"`

	Id       string `bun:",type:uuid,pk,nullzero"`
	Filename string `bun:",notnull,unique,nullzero"`
	FileType string `bun:",notnull,nullzero"`
	Size     uint64 `bun:",notnull,nullzero"`
	IsPublic bool   `bun:",notnull,nullzero"`

	StorageProvider uint8  `bun:",type:uuid,notnull,nullzero"`
	StoragePath     string `bun:",notnull"` // Relative

	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
}

func (m *FileMetadata) ToDomain() domain.FileMetadata {
	return domain.FileMetadata{
		Id:           wrapper.DropError(types.IdFromString(m.Id)),
		Name:         m.Filename,
		Type:         m.FileType,
		Size:         m.Size,
		IsPublic:     m.IsPublic,
		Provider:     domain.StorageProvider(m.StorageProvider),
		ProviderPath: m.StoragePath,
		CreatedAt:    m.CreatedAt,
		LastModified: m.UpdatedAt,
	}
}
