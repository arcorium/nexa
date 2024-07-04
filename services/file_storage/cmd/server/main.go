package main

import (
  "log"
  "nexa/services/file_storage/config"
  sharedConf "nexa/shared/config"
  "nexa/shared/env"
  "nexa/shared/logger"
)

func main() {
  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }

  err := env.LoadEnvs(envName)
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
