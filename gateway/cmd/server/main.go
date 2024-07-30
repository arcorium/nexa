package main

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/env"
  "github.com/arcorium/nexa/shared/logger"
  "github.com/siderolabs/grpc-proxy/proxy"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "log"
  "net"
  "nexa/gateway/config"
  "nexa/gateway/internal/factory"
  "nexa/gateway/internal/handler"
)

func main() {
  var err error
  if config.IsDebug() {
    err = env.LoadEnvs("dev.env")
  } else {
    err = env.LoadEnvs()
  }

  // Init global logger
  logg, err := logger.NewZapLogger(config.IsDebug())
  if err != nil {
    log.Fatalln(err)
  }
  logger.SetGlobal(logg)

  conf, err := sharedConf.Load[config.Server]()
  if err != nil {
    env.LogError(err, -1)
  }

  cred := insecure.NewCredentials()
  client, err := factory.NewClient(&conf.Service,
    grpc.WithTransportCredentials(cred),
    grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.Codec())),
  )
  if err != nil {
    log.Fatalln(err)
  }

  backend := factory.NewBackend(client)

  server := grpc.NewServer(
    grpc.ForceServerCodec(proxy.Codec()),
    grpc.UnknownServiceHandler(
      proxy.TransparentHandler(handler.Director(backend)),
    ),
  )

  logger.Infof("Server listening on %s", conf.Address())
  listen, err := net.Listen("tcp", conf.Address())
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  err = server.Serve(listen)
  log.Fatal(err)
}
