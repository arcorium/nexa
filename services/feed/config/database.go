package config

type Database struct {
  RedisAddress  string `env:"REDIS_ADDRESS"`
  RedisUsername string `env:"REDIS_USERNAME"`
  RedisPassword string `env:"REDIS_PASSWORD"`
}
