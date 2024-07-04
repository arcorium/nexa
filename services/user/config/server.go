package config

import (
  sharedConf "nexa/shared/config"
)

type Server struct {
  sharedConf.Server
  AuthenticationServiceAddress string `env:"AUTHENTICATION_SERVICE_ADDRESS"`
  FileStorageServiceAddress    string `env:"FILE_STORAGE_SERVICE_ADDR"`
  MailerServiceAddress         string `env:"MAILER_SERVICE_ADDR"`
}
