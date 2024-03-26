package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
)

func LoadServer() (*ServerConfig, error) {
	config := &ServerConfig{}
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type ServerConfig struct {
	Ip   string `env:"SERVER_IP"`
	Port uint16 `env:"SERVER_PORT"`
}

func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}
