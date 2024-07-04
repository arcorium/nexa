package main

import (
  "log"
  "nexa/services/authorization/config"
  sharedConf "nexa/shared/config"
  "nexa/shared/env"
  "nexa/shared/logger"
)

func main() {
  err := env.LoadEnvs("dev.env")
  if err != nil {
    log.Fatalln(err)
  }

  // Config
  dbConfig, err := sharedConf.Load[config.Database]()
  if err != nil {
    env.LogError(err, -1)
  }

  serverConfig, err := sharedConf.Load[config.Server]()
  if err != nil {
    log.Fatalln(err)
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
  log.Fatalln(server.Run())
}
