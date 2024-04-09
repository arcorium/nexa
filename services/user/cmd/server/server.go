package main

import (
  "github.com/uptrace/bun"
  "go.uber.org/zap"
  "google.golang.org/grpc"
  "log"
  "net"
  "nexa/services/user/config"
  "nexa/services/user/internal/api/grpc/handler"
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

  grpcServer *grpc.Server
  logger     *zap.Logger

  wg sync.WaitGroup
}

func (s *Server) validationSetup() {
  wrapper.RegisterDefaultNullableValidations(util.GetValidator())
}

func (s *Server) grpcServerSetup() {
  //var err error
  //s.logger, err = zap.NewDevelopment()
  //if err != nil {
  //  log.Fatalln(err)
  //}
  //zapLogger := interceptor.ZapLogger(s.logger)
  //
  //s.server = grpc.NewServer(
  //  //grpc.StatsHandler(otelgrpc.NewServerHandler()),
  //  grpc.ChainUnaryInterceptor(
  //    logging.UnaryServerInterceptor(zapLogger),
  //    otelgrpc.UnaryServerInterceptor(),
  //    recovery.UnaryServerInterceptor(),
  //  ),
  //  grpc.ChainStreamInterceptor(
  //    logging.StreamServerInterceptor(zapLogger),
  //    otelgrpc.StreamServerInterceptor(),
  //    recovery.StreamServerInterceptor(),
  //  ),
  //)
  s.grpcServer = grpc.NewServer()
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
  userHandler.Register(s.grpcServer)

  profileHandler := handler.NewProfileHandler(profileService)
  profileHandler.Register(s.grpcServer)

  return nil
}

func (s *Server) shutdown() {
  log.Println("Server Stopped!")
  s.grpcServer.GracefulStop()
  s.wg.Wait()

  if err := s.db.Close(); err != nil {
    log.Println(err)
  }

  if err := s.logger.Sync(); err != nil {
    log.Println(err)
  }
}

func (s *Server) Run() error {
  listener, err := net.Listen("tcp", s.serverConfig.Address())
  if err != nil {
    return err
  }

  go func() {
    s.wg.Add(1) // NOTE: Needed?
    defer s.wg.Done()

    log.Println("Server Listening on ", s.serverConfig.Address())
    err = s.grpcServer.Serve(listener)
  }()

  quitChan := make(chan os.Signal, 1)
  defer close(quitChan)

  signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
  <-quitChan

  s.shutdown()
  return err
}
