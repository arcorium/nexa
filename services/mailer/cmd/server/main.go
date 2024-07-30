package main

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/env"
  "github.com/arcorium/nexa/shared/logger"
  "log"
  "nexa/services/mailer/config"
)

func main() {
  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }
  _ = env.LoadEnvs(envName)

  // Init global logger
  logg, err := logger.NewZapLogger(config.IsDebug())
  if err != nil {
    log.Fatalln(err)
  }
  logger.SetGlobal(logg)

  // Config
  dbConfig, err := sharedConf.Load[config.Database]()
  if err != nil {
    env.LogError(err, -1)
  }
  logger.Infof("Database Config: %v", dbConfig)

  serverConfig, err := sharedConf.Load[config.Server]()
  if err != nil {
    env.LogError(err, -1)
  }
  logger.Infof("Server Config: %v", serverConfig)

  server, err := NewServer(dbConfig, serverConfig)
  if err != nil {
    log.Fatalln(err)
  }

  if err = server.Run(); err != nil {
    log.Fatalln(err)
  }
}
