package config

import (
  "nexa/services/feed/constant"
  "os"
  "strconv"
  "sync"
)

var once sync.Once
var isDebug bool

// IsDebug return true when this app is on debug or development
// set env RELEASE make this function return false
func IsDebug() bool {
  once.Do(func() {
    s, ok := os.LookupEnv(constant.SERVICE_RELEASE_ENV)
    if !ok {
      isDebug = true
      return
    }

    temp, _ := strconv.ParseBool(s)
    isDebug = !temp
  })
  return isDebug
}
