package main

import (
	"log"
	"nexa/services/file_storage/config"
	sharedConf "nexa/shared/config"
	"nexa/shared/env"
)

func main() {
	err := env.LoadEnvs("dev.env")
	if err != nil {
		log.Fatalln(err)
	}

	// Config
	dbConfig, err := sharedConf.LoadDatabase()
	if err != nil {
		env.LogError(err, -1)
	}

	serverConfig, err := sharedConf.LoadServer[config.Server]()
	if err != nil {
		env.LogError(err, -1)
	}

	server, err := NewServer(dbConfig, serverConfig)
	if err != nil {
		log.Fatalln(err)
	}

	if err = server.Run(); err != nil {
		log.Fatalln(err)
	}
}
