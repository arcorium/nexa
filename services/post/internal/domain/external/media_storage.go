package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IMediaStoreClient interface {
  GetUrls(ctx context.Context, fileIds ...types.Id) ([]string, error)
}
