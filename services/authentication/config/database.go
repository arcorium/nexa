package config

import sharedConf "github.com/arcorium/nexa/shared/config"

type Database struct {
  Postgres sharedConf.PostgresDatabase
  Session  Redis
}

type Redis struct {
  Address  string `env:"REDIS_ADDRESS"`
  Username string `env:"REDIS_USERNAME"`
  Password string `env:"REDIS_PASSWORD"`
}
