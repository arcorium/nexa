package util

import (
  "database/sql"
)

func IsTimeNull(time *sql.NullTime) bool {
  return time.Valid && time.Time.IsZero()
}

func IsStringNull(str *string) bool {
  return str != nil && len(*str) == 0
}
