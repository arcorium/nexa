package config

import (
  "fmt"
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
  MetricPort          uint16 `env:"METRIC_PORT,notEmpty"`
  GrpcExporterAddress string `env:"OTLP_GRPC_EXPORTER_ADDRESS,notEmpty"`
}

func (s *ServerConfig) MetricAddress() string {
  return fmt.Sprintf("%s:%d", s.Ip, s.MetricPort)
}
