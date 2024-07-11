package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
)

type Server struct {
  sharedConf.Server
  PublicKeyPath string `env:"PUBLIC_KEY_PATH"`
  SMTP          SMTP
}

type SMTP struct {
  Host     string `env:"SMTP_HOST,notEmpty"`
  Port     uint16 `env:"SMTP_PORT,notEmpty"`
  Username string `env:"SMTP_USERNAME,notEmpty"`
  Password string `env:"SMTP_PASSWORD,notEmpty"`
}
