package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/mailer/internal/domain/dto"
)

type IFileStorageClient interface {
  GetFiles(ctx context.Context, fileIds ...types.Id) ([]dto.FileAttachment, error)
}
