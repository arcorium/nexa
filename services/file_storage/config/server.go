package config

import "github.com/arcorium/nexa/shared/config"

type Server struct {
  config.Server
  PublicKeyPath string `env:"PUBLIC_KEY_PATH"`
  Storage       Storage
}

type Storage struct {
  Address   string `env:"MINIO_ADDRESS,notEmpty"`
  AccessKey string `env:"MINIO_ACCESS_KEY_ID,notEmpty"`
  SecretKey string `env:"MINIO_SECRET_KEY,notEmpty"`
}
