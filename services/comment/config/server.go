package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "time"
)

type Server struct {
  sharedConf.Server
  CircuitBreaker CircuitBreaker
}

type CircuitBreaker struct {
  MaxRequest       uint32        `env:"HALF_STATE_MAX_REQUEST" envDefault:"5"`
  ResetInterval    time.Duration `env:"HALF_STATE_RESET_INTERVAL" envDefault:"60s"`
  OpenStateTimeout time.Duration `env:"OPEN_STATE_TIMEOUT" envDefault:"30s"`
}
