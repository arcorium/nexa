package entity

import (
	"nexa/shared/types"
	"time"
)

type FileMetadata struct {
	Id       types.Id
	Name     string
	Type     string
	Size     uint64
	IsPublic bool

	Provider     StorageProvider
	ProviderPath string // File path on provider (relative)
	FullPath     string

	CreatedAt    time.Time
	LastModified time.Time
}
