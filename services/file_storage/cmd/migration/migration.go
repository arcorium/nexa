package main

import (
	"log"
	"nexa/services/file_storage/internal/infra/repository/model"
	sharedConf "nexa/shared/config"
	"nexa/shared/database"
	"nexa/shared/env"
)

func main() {
	if err := env.LoadEnvs("dev.env"); err != nil {
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
