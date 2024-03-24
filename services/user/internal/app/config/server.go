package config

type ServerConfig struct {
	Ip   string `env:"SERVER_IP"`
	Port uint16 `env:"SERVER_PORT"`
}
