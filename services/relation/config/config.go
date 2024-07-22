package config

import (
  "nexa/services/relation/constant"
  "os"
  "sync"
)

var once sync.Once
var val bool

// IsDebug return true when this app is on debug or development
// set env RELEASE make this function return false
func IsDebug() bool {
  once.Do(func() {
    _, ok := os.LookupEnv(constant.SERVICE_RELEASE_ENV)
    val = !ok
  })
  return val
}
