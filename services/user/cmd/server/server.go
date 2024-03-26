package main

import (
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
	"os"
	"os/signal"
	"syscall"
)

func NewServer(dbConfig *database.Config, serverConfig *config.ServerConfig) Server {
	return Server{
		dbConfig:     dbConfig,
		serverConfig: serverConfig,
		server:       grpc.NewServer(),
	}
}

type Server struct {
	dbConfig     *database.Config
	serverConfig *config.ServerConfig

	server *grpc.Server
}

func (s *Server) Init() error {

	// Database
	db, err := database.OpenPostgres(s.dbConfig, true, util.Nil[model.User](), util.Nil[model.Profile]())
	if err != nil {
		return err
	}

	// UOW
	userUOW := uow.NewUserUOW(db)

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
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.serverConfig.Address())
	if err != nil {
		return err
	}

	go func() {
		err = s.server.Serve(listener)
	}()

	quitChan := make(chan os.Signal, 1)
	defer close(quitChan)

	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitChan

	s.shutdown()
	return err
}
