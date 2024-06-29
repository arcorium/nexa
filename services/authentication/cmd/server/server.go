package main

import (
  "github.com/uptrace/bun"
  "google.golang.org/grpc"
  "log"
  "nexa/services/authentication/config"
  "nexa/services/authentication/util"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
  "reflect"
  "sync"
)

func NewServer(config *config.Database, serverConfig *config.Server) (*Server, error) {
  svr := &Server{
    dbConfig:     nil,
    serverConfig: nil,
    db:           nil,
    server:       nil,
    wg:           sync.WaitGroup{},
  }
  return svr, svr.Init()
}

type Server struct {
  dbConfig     *config.Database
  serverConfig *config.Server
  db           *bun.DB

  server *grpc.Server

  wg sync.WaitGroup
}

func (s *Server) validationSetup() {
  validator := sharedUtil.GetValidator()
  wrapper.RegisterDefaultNullableValidations(validator)
  util.RegisterValidationTags(validator)
}

func (s *Server) grpcServerSetup() {
  // TODO: Add interceptor
  s.server = grpc.NewServer()
}

func (s *Server) Init() error {
  s.validationSetup()
  s.grpcServerSetup()

}

// serviceDiscovery handle external services discovery
func (s *Server) serviceDiscovery() {
  s.wg.Add(1)
  defer s.wg.Done()

  // TODO: Handle how repository or external client should be renewed which will be used by services
  reflect.TypeFor[Server]().Name()
}

func (s *Server) shutdown() {
  log.Println("Shutting down...")
  s.server.GracefulStop()
  if err := s.db.Close(); err != nil {
    log.Println("Error closing database:", err)
  }
  s.wg.Wait()
}

func (s *Server) Run() {
  go s.serviceDiscovery()
}
