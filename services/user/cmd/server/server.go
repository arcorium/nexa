package main

import (
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"log"
	"net"
	"nexa/services/user/internal/api/grpc/handler"
	"nexa/services/user/internal/app/config"
	"nexa/services/user/internal/app/service"
	"nexa/services/user/internal/app/uow"
	"nexa/services/user/internal/infra/model"
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

	err := svr.Init()
	return svr, err
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
	model.RegisterBunModels(s.db)
	if err != nil {
		return err
	}

	// UOW
	userUOW := uow.NewUserUOW(s.db)

	// Service
	userService := service.NewUser(userUOW)
	profileService := service.NewProfile(userUOW)

	// Api
	userHandler := handler.NewUserHandler(userService)
	userHandler.Register(s.server)

	profileHandler := handler.NewProfileHandler(profileService)
	profileHandler.Register(s.server)

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
