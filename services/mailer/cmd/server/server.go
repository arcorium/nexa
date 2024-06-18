package main

import (
  "context"
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
  "go.uber.org/zap"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
  "log"
  "net"
  "net/http"
  "nexa/services/mailer/config"
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/app/service"
  "nexa/services/mailer/internal/app/uow/pg"
  "nexa/services/mailer/internal/infra/external/mail"
  "nexa/services/mailer/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  sharedInterceptor "nexa/shared/grpc/interceptor"
  "nexa/shared/util"
  "nexa/shared/wrapper"
  "os"
  "os/signal"
  "sync"
  "syscall"
)

func NewServer(dbConfig *sharedConf.Database, serverConfig *config.Server) (*Server, error) {
  svr := &Server{
    dbConfig:     dbConfig,
    serverConfig: serverConfig,
  }

  err := svr.setup()
  return svr, err
}

type Server struct {
  dbConfig     *sharedConf.Database
  serverConfig *config.Server
  db           *bun.DB

  grpcServer     *grpc.Server
  metricServer   *http.Server
  logger         *zap.Logger
  exporter       sdktrace.SpanExporter
  tracerProvider *sdktrace.TracerProvider

  wg sync.WaitGroup
}

func (s *Server) validationSetup() {
  wrapper.RegisterDefaultNullableValidations(util.GetValidator())
}

func (s *Server) grpcServerSetup() error {
  var err error
  s.logger, err = zap.NewDevelopment()
  if err != nil {
    log.Fatalln(err)
  }
  // Log
  zapLogger := sharedInterceptor.ZapLogger(s.logger)

  // Metrics
  metrics := promProv.NewServerMetrics(
    promProv.WithServerHandlingTimeHistogram(
      promProv.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
    ),
  )

  reg := prometheus.NewRegistry()
  reg.MustRegister(metrics)
  exemplarFromCtx := func(ctx context.Context) prometheus.Labels {
    if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
      return prometheus.Labels{"traceID": span.TraceID().String()}
    }
    return nil
  }

  // Trace
  // Exporter
  s.exporter, err = otlptracegrpc.New(context.Background(),
    otlptracegrpc.WithInsecure(),
    otlptracegrpc.WithEndpoint(s.serverConfig.OTLPGRPCCollectorAddress),
  )
  if err != nil {
    return err
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

  s.grpcServer = grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()), // tracing
    grpc.ChainUnaryInterceptor(
      recovery.UnaryServerInterceptor(),
      logging.UnaryServerInterceptor(zapLogger), // logging
      metrics.UnaryServerInterceptor(promProv.WithExemplarFromContext(exemplarFromCtx)),
    ),
    grpc.ChainStreamInterceptor(
      recovery.StreamServerInterceptor(),
      logging.StreamServerInterceptor(zapLogger), // logging
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
  // Database
  var err error
  s.db, err = database.OpenPostgres(s.dbConfig, true)
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

  uow := pg.NewMailUOW(s.db)
  repos := uow.Repositories()
  mailRepo := repos.Mail()
  tagRepo := repos.Tag()

  mailer := mail.NewSMTP(mailRepo, &mail.SMTPConfig{
    Host:     s.serverConfig.SMTPHost,
    Port:     s.serverConfig.SMTPPort,
    Username: s.serverConfig.SMTPUsername,
    Password: s.serverConfig.SMTPPassword,
  })

  mailService := service.NewMail(uow, mailer)
  tagService := service.NewTag(tagRepo)

  return nil
}

func (s *Server) shutdown() {
  s.grpcServer.GracefulStop()
  s.metricServer.Shutdown(context.Background())
  s.wg.Wait()

  // GRPC thingies
  s.exporter.Shutdown(context.Background())

  if err := s.db.Close(); err != nil {
    log.Println(err)
  }

  if err := s.logger.Sync(); err != nil {
    log.Println(err)
  }

  if err := s.tracerProvider.Shutdown(context.Background()); err != nil {
    log.Println(err)
  }

  log.Println("Server Stopped!")
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

    log.Println("Server Listening on ", s.serverConfig.Address())
    err = s.grpcServer.Serve(listener)
    log.Println("Server Stopping ")
    if err != nil {
      log.Println("Server failed to serve:", err)
    }
  }()

  go func() {
    s.wg.Add(1)
    defer s.wg.Done()
    log.Println("Metrics Server Listening on ", s.serverConfig.MetricAddress())

    err = s.metricServer.ListenAndServe()
    log.Println("Metrics Server Stopping")
    if err != nil {
      log.Println("Metrics server failed to serve:", err)
    }
  }()

  quitChan := make(chan os.Signal, 1)
  defer close(quitChan)

  signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
  <-quitChan

  s.shutdown()
  return err
}
