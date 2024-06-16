package config

import (
  "github.com/caarlos0/env/v10"
  "nexa/shared/server"
)

func LoadServer() (*ServerConfig, error) {
  config := &ServerConfig{}
  err := env.Parse(config)
  if err != nil {
    return nil, err
  }

  return config, nil
}

type ServerConfig struct {
  server.Config
}
