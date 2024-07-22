package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IMediaStore interface {
  GetUrls(ctx context.Context, fileIds ...types.Id) ([]string, error)
}
