package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/golang-jwt/jwt/v5"
  "time"
)

type Server struct {
  sharedConf.Server
  JWTAccessTokenExpiration  time.Duration `env:"JWT_ACCESS_TOKEN_EXP" envDefault:"5m"`
  JWTRefreshTokenExpiration time.Duration `env:"JWT_REFRESH_TOKEN_EXP" envDefault:"720h"`
  JWTSigningMethod          string        `env:"JWT_SIGNING_METHOD"`
  PrivateKeyPath            string        `env:"PRIVATE_KEY_PATH"`
  PublicKeyPath             string        `env:"PUBLIC_KEY_PATH"`

  CircuitBreaker CircuitBreaker
  Service        Service

  signingMethod jwt.SigningMethod
}

func (s *Server) SigningMethod() jwt.SigningMethod {
  if s.signingMethod == nil {
    if len(s.JWTSigningMethod) == 0 {
      s.signingMethod = sharedJwt.DefaultSigningMethod
    } else {
      s.signingMethod = jwt.GetSigningMethod(s.JWTSigningMethod)
    }
  }
  return s.signingMethod
}

type Service struct {
  Authorization string `env:"AUTHZ_SERVICE_ADDR,notEmpty"`
  Token         string `env:"TOKEN_SERVICE_ADDR,notEmpty"`
  FileStorage   string `env:"FILE_STORAGE_SERVICE_ADDR,notEmpty"`
  Mailer        string `env:"MAILER_SERVICE_ADDR,notEmpty"`
}

type CircuitBreaker struct {
  MaxRequest       uint32        `env:"HALF_STATE_MAX_REQUEST" envDefault:"5"`
  ResetInterval    time.Duration `env:"HALF_STATE_RESET_INTERVAL" envDefault:"60s"`
  OpenStateTimeout time.Duration `env:"OPEN_STATE_TIMEOUT" envDefault:"30s"`
}
