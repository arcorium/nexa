package main

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/env"
  "log"
  "nexa/services/file_storage/config"
  "nexa/services/file_storage/internal/infra/repository/model"
)

func main() {
  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }
  _ = env.LoadEnvs(envName)

  dbConfig, err := sharedConf.Load[sharedConf.PostgresDatabase]()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgresWithConfig(dbConfig, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  if err = model.CreateTables(db); err != nil {
    log.Fatalln(err)
  }

  log.Println("Succeed migrate database: ", dbConfig.DSN())
}
