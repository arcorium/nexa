package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/user/internal/domain/dto"
)

type IFileStorageClient interface {
  // UploadProfileImage upload file image as public
  UploadProfileImage(ctx context.Context, dto *dto.UploadImageDTO) (types.Id, types.FilePath, error)
  DeleteProfileImage(ctx context.Context, id types.Id) error
}
