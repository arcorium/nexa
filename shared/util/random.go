package util

import (
  "math/rand"
  rand2 "math/rand/v2"
  "strings"
)

func RandomString(length uint64) string {
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

  builder := strings.Builder{}
  builder.Grow(int(length))

  for i := 0; i < int(length); i++ {
    builder.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
  }

  return builder.String()
}

func RandomBool() bool {
  return rand2.IntN(1) == 1
}
