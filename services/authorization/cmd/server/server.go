package main

import (
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"log"
	"net"
	"nexa/services/authorization/internal/api/grpc/handler"
	"nexa/services/authorization/internal/app/config"
	"nexa/services/authorization/internal/app/service"
	"nexa/services/authorization/internal/infra/repository/factory"
	"nexa/shared/database"
	"nexa/shared/util"
	"nexa/shared/wrapper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func NewServer(dbConfig *database.Config, serverConfig *config.ServerConfig) (*Server, error) {
	svr := &Server{
		dbConfig:     dbConfig,
		serverConfig: serverConfig,
	}
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
	s.server = grpc.NewServer()
}

func (s *Server) Init() error {
	s.validationSetup()
	s.grpcServerSetup()

	// Database
	var err error
	s.db, err = database.OpenPostgres(s.dbConfig, true)
	if err != nil {
		return err
	}

	// Repositories
	repositories := factory.NewPostgresRepositories(s.db)

	// Services
	actionSvc := service.NewAction(repositories.Action)
	permissionSvc := service.NewPermission(repositories.Permission)
	resourceSvc := service.NewResource(repositories.Resource)
	roleSvc := service.NewRole(repositories.Role)

	// Api
	actionApi := handler.NewAction(actionSvc)
	actionApi.Register(s.server)
	permissionApi := handler.NewPermission(permissionSvc)
	permissionApi.Register(s.server)
	resourceApi := handler.NewResource(resourceSvc)
	resourceApi.Register(s.server)
	roleApi := handler.NewRole(roleSvc)
	roleApi.Register(s.server)

	return nil
}

func (s *Server) shutdown() {
	log.Println("Server Stopped!")
	s.server.GracefulStop()
	if err := s.db.Close(); err != nil {
		log.Println(err)
	}
	s.wg.Wait()
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.serverConfig.Address())
	if err != nil {
		return err
	}

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		log.Println("Server Listening on ", s.serverConfig.Address())
		err = s.server.Serve(listener)
	}()

	quitChan := make(chan os.Signal, 1)
	defer close(quitChan)

	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitChan

	s.shutdown()
	return err
}
