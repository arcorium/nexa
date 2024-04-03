package server

import "fmt"

type Config struct {
	Ip   string `env:"SERVER_IP" envDefault:"localhost"`
	Port uint16 `env:"SERVER_PORT,notEmpty"`
}

func (s *Config) Address() string {
	return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}
