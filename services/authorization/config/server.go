package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
)

type Server struct {
  sharedConf.Server
  PublicKeyPath string `env:"PUBLIC_KEY_PATH"`
}
