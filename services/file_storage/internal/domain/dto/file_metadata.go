package dto

import (
  "nexa/shared/wrapper"
  "time"
)

type UpdateFileMetadataDTO struct {
  Id       string `json:"id" validate:"required,uuid4"`
  Name     wrapper.Nullable[string]
  IsPublic wrapper.Nullable[bool]
}

type FileMetadataResponseDTO struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Size uint64 `json:"size"`
  Path string `json:"path"`

  CreatedAt    time.Time `json:"created_at"`
  LastModified time.Time `json:"last_modified"`
}
