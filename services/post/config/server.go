package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "time"
)

type Server struct {
  sharedConf.Server
  PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
  CircuitBreaker CircuitBreaker
  Service        Service
}

type CircuitBreaker struct {
  MaxRequest       uint32        `env:"HALF_STATE_MAX_REQUEST" envDefault:"5"`
  ResetInterval    time.Duration `env:"HALF_STATE_RESET_INTERVAL" envDefault:"60s"`
  OpenStateTimeout time.Duration `env:"OPEN_STATE_TIMEOUT" envDefault:"30s"`
}

type Service struct {
  Comment      string `env:"COMMENT_SERVICE_ADDRESS"`
  Follow       string `env:"FOLLOW_SERVICE_ADDRESS"`
  MediaStorage string `env:"MEDIA_STORAGE_SERVICE_ADDRESS"`
  Reaction     string `env:"REACTION_SERVICE_ADDRESS"`
  User         string `env:"USER_SERVICE_ADDRESS"`
}
