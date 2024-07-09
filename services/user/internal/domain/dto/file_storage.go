package dto

import "github.com/arcorium/nexa/shared/types"

type UploadImageDTO struct {
  Filename string `json:"filename" validate:"required"`
  Data     []byte `json:"data" validate:"required"`
}

type UpdateImageDTO struct {
  Id       types.Id `json:"id" validate:"required,uuid4"`
  Filename string   `json:"filename" validate:"required"`
  Data     []byte   `json:"data" validate:"required"`
}
