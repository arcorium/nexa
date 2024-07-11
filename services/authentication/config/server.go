package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/golang-jwt/jwt/v5"
  "time"
)

type Server struct {
  sharedConf.Server
  TokenExpiration           time.Duration `env:"TOKEN_EXPIRATION" envDefault:"24h"`
  JWTAccessTokenExpiration  time.Duration `env:"JWT_ACCESS_TOKEN_EXP" envDefault:"5m"`
  JWTRefreshTokenExpiration time.Duration `env:"JWT_REFRESH_TOKEN_EXP" envDefault:"30d"`
  JWTSigningMethod          string        `env:"JWT_SIGNING_METHOD"`
  PrivateKeyPath            string        `env:"PRIVATE_KEY_PATH"`
  PublicKeyPath             string        `env:"PUBLIC_KEY_PATH"`

  Service Service

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
  User          string `env:"USER_SERVICE_ADDR,notEmpty"`
  Mailer        string `env:"MAILER_SERVICE_ADDR,notEmpty"`
}
