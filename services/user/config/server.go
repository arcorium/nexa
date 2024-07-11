package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
)

type Server struct {
  sharedConf.Server
  Service       Service
  PublicKeyPath string `env:"PUBLIC_KEY_PATH"`
}

type Service struct {
  Authentication string `env:"AUTHN_SERVICE_ADDR,notEmpty"`
  FileStorage    string `env:"FILE_STORAGE_SERVICE_ADDR,notEmpty"`
  Mailer         string `env:"MAILER_SERVICE_ADDR,notEmpty"`
}
