package config

import "nexa/shared/config"

type Server struct {
	config.Server
	MinIOAddress   string `env:"MINIO_ADDRESS,notEmpty"`
	MinIOAccessKey string `env:"MINIO_ACCESS_KEY_ID,notEmpty"`
	MinIOSecretKey string `env:"MINIO_SECRET_KEY,notEmpty"`
}
