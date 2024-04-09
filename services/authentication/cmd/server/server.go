package main

import (
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"log"
	"nexa/services/authentication/internal/app/config"
	"nexa/shared/database"
	"nexa/shared/util"
	"nexa/shared/wrapper"
	"reflect"
	"sync"
)

func NewServer(config *database.Config, serverConfig *config.ServerConfig) (*Server, error) {
	svr := &Server{}
	return svr, svr.Init()
}

type Server struct {
	dbConfig     *database.Config
	serverConfig *config.ServerConfig
	db           *bun.DB

	server *grpc.Server

	wg sync.WaitGroup
}

func (s *Server) validationSetup() {
	wrapper.RegisterDefaultNullableValidations(util.GetValidator())
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
