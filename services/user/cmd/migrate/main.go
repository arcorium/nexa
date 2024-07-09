package main

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "log"
  "nexa/services/user/config"
  "nexa/services/user/internal/infra/repository/model"
)

func main() {
  dbConfig, err := sharedConf.Load[config.PostgresDatabase]()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgres(dbConfig.DSN(), dbConfig.IsSecure, dbConfig.Timeout, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  if err = model.CreateTables(db); err != nil {
    log.Fatalln(err)
  }
}
