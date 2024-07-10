package main

import (
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "log"
  "nexa/services/authorization/internal/infra/repository/model"
)

func main() {
  dbConfig, err := sharedConf.Load[sharedConf.PostgresDatabase]()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgresWithConfig(dbConfig, true)
  //db, err := database.OpenPostgres("postgres://postgres:password@localhost:5432/postgres", false, 0, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  if err = model.CreateTables(db); err != nil {
    log.Fatalln(err)
  }

  log.Println("Success migrate database: ", dbConfig.DSN())
}