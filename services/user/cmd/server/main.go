package main

import (
	"log"
	"nexa/services/user/internal/app/config"
	"nexa/shared/database"
	"nexa/shared/env"
)

func main() {
	err := env.LoadEnvs("dev.env")
	if err != nil {
		log.Fatalln(err)
	}

	// Config
	dbConfig, err := database.LoadConfig()
	if err != nil {
		env.LogError(err, -1)
	}

	serverConfig, err := config.LoadServer()
	if err != nil {
		env.LogError(err, -1)
	}

	server, err := NewServer(dbConfig, serverConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(server.Run())
}
