package util

import (
  "github.com/arcorium/nexa/shared/jwt"
  rand2 "math/rand/v2"
)

const alphanumericString = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const fullString = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandomString(length uint16) string {
  if length == 0 {
    return jwt.GenerateRefreshToken()
  }

  var result []byte
  for range length {
    n := rand2.N(len(fullString))
    result = append(result, fullString[n])
  }
  return string(result)
}

func RandomPIN(alphanumeric bool, length uint16) string {
  var result []byte

  // Set default length
  if length == 0 {
    length = 6
  }

  count := 0xA
  if alphanumeric {
    count = len(alphanumericString)
  }

  for range length {
    n := rand2.N(count)
    result = append(result, alphanumericString[n])
  }
  return string(result)
}
