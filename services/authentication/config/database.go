package config

import sharedConf "github.com/arcorium/nexa/shared/config"

type Database struct {
  sharedConf.PostgresDatabase
  RedisAddress  string `env:"REDIS_ADDRESS"`
  RedisUsername string `env:"REDIS_USERNAME"`
  RedisPassword string `env:"REDIS_PASSWORD"`
}
