package config

import (
  sharedConf "nexa/shared/config"
)

type Server struct {
  sharedConf.Server
  SMTPHost     string `env:"SMTP_HOST,notEmpty"`
  SMTPPort     uint16 `env:"SMTP_PORT,notEmpty"`
  SMTPUsername string `env:"SMTP_USERNAME,notEmpty"`
  SMTPPassword string `env:"SMTP_PASSWORD,notEmpty"`
}
