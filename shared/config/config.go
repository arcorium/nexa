package config

import "github.com/caarlos0/env/v10"

func Load[T any]() (*T, error) {
  var conf T
  err := env.Parse(&conf)
  if err != nil {
    return nil, err
  }

  return &conf, nil
}
