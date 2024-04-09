package main

import (
  "log"
  "nexa/services/user/internal/infra/model"
  "nexa/shared/database"
  "nexa/shared/env"
)

func main() {
  if err := env.LoadEnvs("dev.env"); err != nil {
    log.Println(err)
  }

  dbConfig, err := database.LoadConfig()
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
