package env

import "os"

func GetDefaulted(key string, defaults string) string {
  val := os.Getenv(key)
  if val == "" {
    return defaults
  }
  return val
}

func GetDefaultedP(key string, defaults string) *string {
  val := os.Getenv(key)
  if val == "" {
    return &defaults
  }
  return nil
}
