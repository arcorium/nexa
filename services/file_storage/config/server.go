package config

import "github.com/arcorium/nexa/shared/config"

type Server struct {
  config.Server
  Storage       Storage
  PublicKeyPath string `env:"PUBLIC_KEY_PATH"`
  BucketName    string `env:"BUCKET_NAME,notEmpty"`
}

type Storage struct {
  Address   string `env:"MINIO_ADDRESS,notEmpty"`
  AccessKey string `env:"MINIO_ACCESS_KEY_ID,notEmpty"`
  SecretKey string `env:"MINIO_SECRET_KEY,notEmpty"`
}
