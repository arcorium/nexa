package main

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/env"
  "github.com/arcorium/nexa/shared/logger"
  "log"
  "nexa/services/feed/config"
)

func main() {
  var err error
  if config.IsDebug() {
    err = env.LoadEnvs("dev.env")
  } else {
    err = env.LoadEnvs()
  }

  // Config
  dbConfig, err := sharedConf.Load[config.Database]()
  if err != nil {
    env.LogError(err, -1)
  }

  serverConfig, err := sharedConf.Load[config.Server]()
  if err != nil {
    env.LogError(err, -1)
  }

  // Init global logger
  logg, err := logger.NewZapLogger(config.IsDebug())
  if err != nil {
    log.Fatalln(err)
  }
  logger.SetGlobal(logg)

  server, err := NewServer(dbConfig, serverConfig)
  if err != nil {
    log.Fatalln(err)
  }

  if err = server.Run(); err != nil {
    log.Fatalln(err)
  }
}
