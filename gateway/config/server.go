package config

import "fmt"

type Server struct {
  Ip      string `env:"SERVER_IP" envDefault:"0.0.0.0"`
  Port    int    `env:"SERVER_PORT" envDefault:"8080"`
  Service Service
}

func (s *Server) Address() string {
  return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}

type Service struct {
  Authn    string `env:"AUTHN_SERVICE_ADDRESS,notEmpty"`
  Authz    string `env:"AUTHZ_SERVICE_ADDRESS,notEmpty"`
  Comment  string `env:"COMMENT_SERVICE_ADDRESS,notEmpty"`
  Feed     string `env:"FEED_SERVICE_ADDRESS,notEmpty"`
  Storage  string `env:"STORAGE_SERVICE_ADDRESS,notEmpty"`
  Mailer   string `env:"MAILER_SERVICE_ADDRESS,notEmpty"`
  Post     string `env:"POST_SERVICE_ADDRESS,notEmpty"`
  Reaction string `env:"REACTION_SERVICE_ADDRESS,notEmpty"`
  Relation string `env:"RELATION_SERVICE_ADDRESS,notEmpty"`
  //Token    string `env:"TOKEN_SERVICE_ADDRESS"`
}
