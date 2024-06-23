package external

import (
  "context"
  "nexa/services/user/internal/domain/dto"
  "nexa/shared/types"
)

type IFileStorageClient interface {
  UploadProfileImage(ctx context.Context, dto *dto.UploadImageDTO) (types.Id, error)
  UpdateProfileImage(ctx context.Context, dto *dto.UpdateImageDTO) (types.Id, error)
}
