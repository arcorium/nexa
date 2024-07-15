package config

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "time"
)

type Server struct {
  sharedConf.Server
  TokenExpiration    TokenExpiration
  GeneralTokenExpiry time.Duration `env:"GENERAL_TOKEN_EXPIRATION" envDefault:"0"` // Used for use single expiry time for all kind of token usage
}

func (s *Server) IsGeneralExpiry() bool {
  return s.GeneralTokenExpiry != 0
}

type TokenExpiration struct {
  EmailVerification time.Duration `env:"EMAIL_VERIFICATION_TOKEN_EXPIRY" envDefault:"24h"`
  ResetPassword     time.Duration `env:"RESET_PASSWORD_TOKEN_EXPIRY" envDefault:"24h"`
  Login             time.Duration `env:"LOGIN_TOKEN_EXPIRY" envDefault:"24h"`
  Other             time.Duration `env:"OTHER_TOKEN_EXPIRY" envDefault:"24h"`
}
