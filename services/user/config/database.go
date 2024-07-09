package config

import (
  "fmt"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "sync"
  "time"
)

type Database struct {
  sharedConf.Database
}

type PostgresDatabase struct {
  //Host string `env:"PG_HOST,notEmpty"`
  //Port uint16 `env:"PG_PORT,notEmpty"`
  Address string `env:"PG_ADDRESS,notEmpty"`

  Username  string `env:"PG_USERNAME,notEmpty"`
  Password  string `env:"PG_PASSWORD"`
  Name      string `env:"PG_DB_NAME,notEmpty"`
  Parameter string `env:"PG_PARAMETER"`
  IsSecure  bool   `env:"PG_IS_SECURE" envDefault:"false"`

  Timeout time.Duration `env:"PG_TIMEOUT" envDefault:"60s"`

  dsn     string
  dsnOnce sync.Once
}

func (c *PostgresDatabase) DSN() string {
  c.dsnOnce.Do(func() {
    password := ""
    if len(c.Password) != 0 {
      password = ":" + c.Password
    }
    param := ""
    if len(c.Parameter) != 0 {
      param = "?" + c.Parameter
    }

    c.dsn = fmt.Sprintf("postgres://%s%s@%s/%s%s", c.Username, password,
      c.Address, c.Name, param)
  })

  return c.dsn
}
