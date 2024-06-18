package main

import (
  "log"
  "nexa/services/file_storage/config"
  "nexa/services/file_storage/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/env"
)

func main() {
  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }

  if err := env.LoadEnvs(envName); err != nil {
    log.Println(err)
  }

  dbConfig, err := sharedConf.LoadDatabase()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgres(dbConfig, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  if err = model.CreateTables(db); err != nil {
    log.Fatalln(err)
  }
}
