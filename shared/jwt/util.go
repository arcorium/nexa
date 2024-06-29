package jwt

import (
  "crypto/sha512"
  "encoding/hex"
  "github.com/google/uuid"
)

// GenerateRefreshToken does hash the uuid with SHA512 to generate random string
func GenerateRefreshToken() string {
  uuids := uuid.NewString()
  hash := sha512.New()
  hash.Write([]byte(uuids))
  return hex.EncodeToString(hash.Sum(nil))
}
