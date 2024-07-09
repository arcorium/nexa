package optional

import "github.com/arcorium/nexa/shared/types"

var (
  NullString = Null[string]()
  NullId     = Null[types.Id]()
)
