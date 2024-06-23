package config

import (
  "github.com/golang-jwt/jwt/v5"
  "nexa/shared/server"
  "time"
)

type Server struct {
  server.Config
  TokenExpiration    time.Duration
  JWTTokenExpiration time.Duration
  JWTSigningMethod   string
  JWTSecretKey       string

  UserServiceName string

  signingMethod jwt.SigningMethod
}

func (s *Server) SigningMethod() jwt.SigningMethod {
  if s.signingMethod == nil {
    s.signingMethod = jwt.GetSigningMethod(s.JWTSigningMethod)
  }
  return s.signingMethod
}

func (s *Server) SecretKey() []byte {
  return []byte(s.JWTSecretKey)
}

func (s *Server) KeyFunc() jwt.Keyfunc {
  return func(token *jwt.Token) (interface{}, error) {
    return s.SecretKey(), nil
  }
}
