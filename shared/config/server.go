package config

import "fmt"

type Server struct {
	Ip   string `env:"SERVER_IP" envDefault:"localhost"`
	Port uint16 `env:"SERVER_PORT,notEmpty"`

	MetricPort          uint16 `env:"METRIC_PORT,notEmpty"`
	GrpcExporterAddress string `env:"OTLP_GRPC_EXPORTER_ADDRESS,notEmpty"`
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}

func (s *Server) MetricAddress() string {
	return fmt.Sprintf("%s:%d", s.Ip, s.MetricPort)
}
