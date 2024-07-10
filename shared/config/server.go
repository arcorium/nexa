package config

import "fmt"

type Server struct {
  Ip   string `env:"SERVER_IP" envDefault:"0.0.0.0"`
  Port uint16 `env:"SERVER_PORT" envDefault:"8080"`

  MetricPort               uint16 `env:"METRIC_PORT" envDefault:"8081"`
  OTLPGRPCCollectorAddress string `env:"OTLP_GRPC_COLLECTOR_ADDRESS,notEmpty"`
}

func (s *Server) Address() string {
  return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}

func (s *Server) MetricAddress() string {
  return fmt.Sprintf("%s:%d", s.Ip, s.MetricPort)
}
