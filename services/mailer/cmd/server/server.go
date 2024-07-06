package main

import (
  "context"
  "errors"
  promProv "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/collectors"
  "github.com/prometheus/client_golang/prometheus/promhttp"
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
  "google.golang.org/grpc/reflection"
  "net"
  "net/http"
  "nexa/services/mailer/config"
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/api/grpc/handler"
  inter "nexa/services/mailer/internal/api/grpc/interceptor"
  "nexa/services/mailer/internal/app/service"
  "nexa/services/mailer/internal/app/uow/pg"
  external2 "nexa/services/mailer/internal/domain/external"
  "nexa/services/mailer/internal/infra/external"
  "nexa/services/mailer/internal/infra/repository/model"
  "nexa/shared/database"
  "nexa/shared/grpc/interceptor"
  "nexa/shared/logger"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
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
  dbConfig     *config.Database
  serverConfig *config.Server
  db           *bun.DB
  mailClient   external2.IMail

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

  s.grpcServer = grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()), // tracing
    grpc.ChainUnaryInterceptor(
      recovery.UnaryServerInterceptor(),
      logging.UnaryServerInterceptor(zapLogger), // logging
      interceptor.UnaryServerAuth(inter.Auth, inter.AuthSelector),
      metrics.UnaryServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
    grpc.ChainStreamInterceptor(
      recovery.StreamServerInterceptor(),
      logging.StreamServerInterceptor(zapLogger), // logging
      interceptor.StreamServerAuth(inter.Auth, inter.AuthSelector),
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
  s.db, err = database.OpenPostgres(&s.dbConfig.Database, true)
  if err != nil {
    return err
  }
  // Add trace hook
  s.db.AddQueryHook(bunotel.NewQueryHook(
    bunotel.WithFormattedQueries(true),
  ))
  model.RegisterBunModels(s.db)

  return nil
}

func (s *Server) setup() error {
  s.validationSetup()

  if err := s.grpcServerSetup(); err != nil {
    return err
  }
  if err := s.databaseSetup(); err != nil {
    return err
  }

  // External client
  {
    smtpClient, err := external.NewSMTP(&external.SMTPConfig{
      Host:     s.serverConfig.SMTPHost,
      Port:     s.serverConfig.SMTPPort,
      Username: s.serverConfig.SMTPUsername,
      Password: s.serverConfig.SMTPPassword,
    })
    if err != nil {
      return err
    }
    s.mailClient = smtpClient
  }

  // Repository
  mailUOW := pg.NewMailUOW(s.db)
  repos := mailUOW.Repositories()

  // Service
  mailService := service.NewMail(mailUOW, s.mailClient)
  tagService := service.NewTag(repos.Tag())

  // GRPC Handler
  mailHandler := handler.NewMail(mailService)
  mailHandler.Register(s.grpcServer)

  tagHandler := handler.NewTag(tagService)
  tagHandler.Register(s.grpcServer)

  return nil
}

func (s *Server) shutdown() {
  ctx := context.Background()

  s.grpcServer.GracefulStop()
  s.metricServer.Shutdown(ctx)
  s.wg.Wait()

  // OTEL
  s.exporter.Shutdown(ctx)

  if err := s.tracerProvider.Shutdown(ctx); err != nil {
    logger.Warn(err.Error())
  }

  // External
  if err := s.mailClient.Close(ctx); err != nil {
    logger.Warn(err.Error())
  }

  // Database
  if err := s.db.Close(); err != nil {
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
