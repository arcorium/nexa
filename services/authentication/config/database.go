package config

import sharedConf "nexa/shared/config"

type Database struct {
  sharedConf.Database
  SessionDatabaseAddress  string `env:"REDIS_ADDRESS"`
  SessionDatabaseUsername string `env:"REDIS_USERNAME"`
  SessionDatabasePassword string `env:"REDIS_PASSWORD"`
}
