package main

import (
  "context"
  "crypto/rsa"
  "errors"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/grpc/interceptor"
  "github.com/arcorium/nexa/shared/grpc/interceptor/authz"
  "github.com/arcorium/nexa/shared/logger"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/golang-jwt/jwt/v5"
  promProv "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/collectors"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/redis/go-redis/v9"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/extra/bunotel"
  "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  "go.opentelemetry.io/otel/propagation"
  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/health"
  "google.golang.org/grpc/health/grpc_health_v1"
  "google.golang.org/grpc/reflection"
  "net"
  "net/http"
  "nexa/services/authentication/config"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/internal/api/grpc/handler"
  inter "nexa/services/authentication/internal/api/grpc/interceptor"
  "nexa/services/authentication/internal/app/service"
  "nexa/services/authentication/internal/infra/external"
  "nexa/services/authentication/internal/infra/repository/model"
  "nexa/services/authentication/internal/infra/repository/pg"
  redisRepo "nexa/services/authentication/internal/infra/repository/redis"
  "nexa/services/authentication/util"
  "os"
  "os/signal"
  "sync"
  "syscall"
)

func NewServer(dbConfig *config.Database, serverConfig *config.Server) (*Server, error) {
  svr := &Server{
    dbConfig:     dbConfig,
    serverConfig: serverConfig,
  }

  err := svr.setup()
  return svr, err
}

type Server struct {
  dbConfig              *config.Database
  serverConfig          *config.Server
  tokenDb               *bun.DB
  credDb                *redis.Client
  grpcClientConnections []*grpc.ClientConn
  publicKey             *rsa.PublicKey
  privateKey            *rsa.PrivateKey

  grpcServer     *grpc.Server
  metricServer   *http.Server
  logger         logger.ILogger
  exporter       sdktrace.SpanExporter
  tracerProvider *sdktrace.TracerProvider

  wg sync.WaitGroup
}

func (s *Server) validationSetup() {
  validator := sharedUtil.GetValidator()
  types.RegisterDefaultNullableValidations(validator)
  util.RegisterValidationTags(validator)
}

func (s *Server) setupOtel() (*promProv.ServerMetrics, *prometheus.Registry, error) {
  var err error
  // Metrics
  metrics := promProv.NewServerMetrics(
    promProv.WithServerHandlingTimeHistogram(
      promProv.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
    ),
  )

  reg := prometheus.NewRegistry()
  reg.MustRegister(metrics)

  // Trace
  // Exporter
  s.exporter, err = otlptracegrpc.New(context.Background(),
    otlptracegrpc.WithInsecure(),
    otlptracegrpc.WithEndpoint(s.serverConfig.OTLPGRPCCollectorAddress),
  )
  if err != nil {
    return nil, nil, err
  }

  // Resource
  res, err := resource.New(context.Background(),
    resource.WithAttributes(
      semconv.ServiceName(constant.SERVICE_NAME),
      semconv.ServiceVersion(constant.SERVICE_VERSION),
    ))

  bsp := sdktrace.NewBatchSpanProcessor(s.exporter)
  s.tracerProvider = sdktrace.NewTracerProvider(
    sdktrace.WithSampler(sdktrace.AlwaysSample()),
    sdktrace.WithSpanProcessor(bsp),
    sdktrace.WithResource(res),
  )

  otel.SetTracerProvider(s.tracerProvider)
  otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
    propagation.TraceContext{}, propagation.Baggage{},
  ))

  return metrics, reg, nil
}

func (s *Server) grpcServerSetup() error {
  // Log
  s.logger = logger.GetGlobal()
  zaplogger, ok := s.logger.(*logger.ZapLogger)
  if !ok {
    return errors.New("logger is not of expected type, expected zap")
  }
  zapLogger := interceptor.ZapLogger(zaplogger.Internal)

  metrics, reg, err := s.setupOtel()
  if err != nil {
    return err
  }

  exemplarFromCtx := func(ctx context.Context) prometheus.Labels {
    if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
      return prometheus.Labels{"traceID": span.TraceID().String()}
    }
    return nil
  }

  authorizationConfig := authz.NewUserConfig(s.publicKey, inter.PermissionCheck)

  s.grpcServer = grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()), // tracing
    grpc.ChainUnaryInterceptor(
      recovery.UnaryServerInterceptor(),
      logging.UnaryServerInterceptor(zapLogger), // logging
      authz.UserUnaryServerInterceptor(&authorizationConfig,
        authz.SkipSelector(inter.AuthSkipSelector)),
      metrics.UnaryServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
    grpc.ChainStreamInterceptor(
      recovery.StreamServerInterceptor(),
      logging.StreamServerInterceptor(zapLogger), // logging
      authz.UserStreamServerInterceptor(&authorizationConfig,
        authz.SkipSelector(inter.AuthSkipSelector)),
      metrics.StreamServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
  )

  if config.IsDebug() {
    reflection.Register(s.grpcServer)
    reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
    reg.MustRegister(collectors.NewGoCollector())
  }

  metrics.InitializeMetrics(s.grpcServer)
  // Metric endpoint
  mux := http.NewServeMux()
  mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
    Registry:          reg,
    EnableOpenMetrics: true,
  }))

  s.metricServer = &http.Server{Handler: mux, Addr: s.serverConfig.MetricAddress()}

  return nil
}

func (s *Server) databaseSetup() error {
  // Token Database
  var err error
  s.tokenDb, err = database.OpenPostgresWithConfig(&s.dbConfig.Postgres, config.IsDebug())
  if err != nil {
    return err
  }
  // Add trace hook
  s.tokenDb.AddQueryHook(bunotel.NewQueryHook(
    bunotel.WithFormattedQueries(true),
  ))
  model.RegisterBunModels(s.tokenDb)

  // Credential Database
  s.credDb = redis.NewClient(&redis.Options{
    Addr:     s.dbConfig.Session.Address,
    Username: s.dbConfig.Session.Username,
    Password: s.dbConfig.Session.Password,
  })

  _, err = s.credDb.Ping(context.Background()).Result()
  if err != nil {
    return err
  }

  return nil
}

func (s *Server) setupKey() error {
  // Get public key from PEM
  pubkeyPath := "pubkey.pem"
  if s.serverConfig.PublicKeyPath != "" {
    pubkeyPath = s.serverConfig.PublicKeyPath
  }
  pub, err := os.ReadFile(pubkeyPath)
  if err != nil {
    return err
  }

  publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
  if err != nil {
    return err
  }
  s.publicKey = publicKey

  // Get private key from PEM
  privkeyPath := "privkey.pem"
  if s.serverConfig.PrivateKeyPath != "" {
    pubkeyPath = s.serverConfig.PrivateKeyPath
  }
  priv, err := os.ReadFile(privkeyPath)
  if err != nil {
    return err
  }

  privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
  if err != nil {
    return err
  }

  s.privateKey = privKey
  return nil
}

func (s *Server) setup() error {
  s.validationSetup()

  if err := s.setupKey(); err != nil {
    return err
  }
  if err := s.grpcServerSetup(); err != nil {
    return err
  }
  if err := s.databaseSetup(); err != nil {
    return err
  }

  // External client
  creds := grpc.WithTransportCredentials(insecure.NewCredentials())
  roleConn, err := grpc.NewClient(s.serverConfig.Service.Authorization, creds)
  if err != nil {
    return err
  }
  roleClient := external.NewRoleClient(roleConn)

  userConn, err := grpc.NewClient(s.serverConfig.Service.User, creds)
  if err != nil {
    return err
  }
  userClient := external.NewUserClient(userConn)

  mailConn, err := grpc.NewClient(s.serverConfig.Service.Mailer, creds)
  if err != nil {
    return err
  }
  mailClient := external.NewMailClient(mailConn)
  s.grpcClientConnections = append(s.grpcClientConnections, roleConn, userConn, mailConn)

  // Repository
  tokenRepo := pg.NewToken(s.tokenDb)
  credRepo := redisRepo.NewCredential(s.credDb, nil) // Use default config

  // Service
  tokenService := service.NewToken(tokenRepo, service.TokenServiceConfig{
    VerificationTokenExpiration: s.serverConfig.TokenExpiration,
    ResetTokenExpiration:        s.serverConfig.TokenExpiration,
  })

  credService := service.NewCredential(credRepo, tokenRepo, roleClient, userClient, mailClient, service.CredentialServiceConfig{
    SigningMethod:          s.serverConfig.SigningMethod(),
    AccessTokenExpiration:  s.serverConfig.JWTAccessTokenExpiration,
    RefreshTokenExpiration: s.serverConfig.JWTRefreshTokenExpiration,
    PrivateKey:             s.privateKey,
    PublicKey:              s.publicKey,
  })

  // GRPC Handler
  tokenHandler := handler.NewToken(tokenService)
  tokenHandler.RegisterHandler(s.grpcServer)

  credHandler := handler.NewCredential(credService)
  credHandler.RegisterHandler(s.grpcServer)

  // Health check
  healthHandler := health.NewServer()
  grpc_health_v1.RegisterHealthServer(s.grpcServer, healthHandler)
  healthHandler.SetServingStatus(constant.SERVICE_NAME, grpc_health_v1.HealthCheckResponse_SERVING)

  return nil
}

func (s *Server) shutdown() {
  s.grpcServer.GracefulStop()
  s.metricServer.Shutdown(context.Background())
  s.wg.Wait()

  // OTEL
  s.exporter.Shutdown(context.Background())

  if err := s.tracerProvider.Shutdown(context.Background()); err != nil {
    logger.Warn(err.Error())
  }

  // External Clients
  for _, conn := range s.grpcClientConnections {
    if err := conn.Close(); err != nil {
      logger.Warn(err.Error())
    }
  }

  // Database
  if err := s.tokenDb.Close(); err != nil {
    logger.Warn(err.Error())
  }

  if err := s.credDb.Close(); err != nil {
    logger.Warn(err.Error())
  }

  logger.Info("Server Stopped!")
}

func (s *Server) Run() error {
  listener, err := net.Listen("tcp", s.serverConfig.Address())
  if err != nil {
    return err
  }

  // Run grpc server
  go func() {
    s.wg.Add(1)
    defer s.wg.Done()

    logger.Infof("Server Listening on %s", s.serverConfig.Address())

    err = s.grpcServer.Serve(listener)
    logger.Info("Server Stopping ")
    if err != nil {
      logger.Warnf("Server failed to serve: %s", err)
    }
  }()

  go func() {
    s.wg.Add(1)
    defer s.wg.Done()
    logger.Infof("Metrics Server Listening on %s", s.serverConfig.MetricAddress())

    err = s.metricServer.ListenAndServe()
    logger.Info("Metrics Server Stopping")
    if err != nil {
      logger.Warnf("Metrics server failed to serve: %s", err)
    }
  }()

  quitChan := make(chan os.Signal, 1)
  defer close(quitChan)

  signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
  <-quitChan

  s.shutdown()
  return err
}
