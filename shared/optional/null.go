package optional

import "nexa/shared/types"

var (
  NullString = Null[string]()
  NullId     = Null[types.Id]()
)
