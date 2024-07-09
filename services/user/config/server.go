package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
)

type Server struct {
  sharedConf.Server
  AuthenticationServiceAddress string `env:"AUTHENTICATION_SERVICE_ADDRESS,notEmpty"`
  FileStorageServiceAddress    string `env:"FILE_STORAGE_SERVICE_ADDR,notEmpty"`
  MailerServiceAddress         string `env:"MAILER_SERVICE_ADDR,notEmpty"`
}
