package dto

type FileResponseDTO struct {
  Name string `json:"name"`
  Size uint64 `json:"size"`
  Data []byte `json:"data"`
}

type FileStoreDTO struct {
  Name     string `json:"name" validate:"required"`
  Data     []byte `json:"data" validate:"required"`
  IsPublic bool   `json:"is_public" validate:"required"`
}
