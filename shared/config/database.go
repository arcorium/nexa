package config

import (
  "fmt"
  "sync"
  "time"
)

func LoadDatabase() (*Database, error) {
  return Load[Database]()
}

type Database struct {
  Protocol string `env:"DB_PROTOCOL,notEmpty"`
  Host     string `env:"DB_HOST,notEmpty"`
  Port     uint16 `env:"DB_PORT,notEmpty"`

  Username  string `env:"DB_USERNAME,notEmpty"`
  Password  string `env:"DB_PASSWORD"`
  Name      string `env:"DB_NAME,notEmpty"`
  Parameter string `env:"DB_PARAMETER"`
  IsSecure  bool   `env:"DB_IS_SECURE" envDefault:"false"`

  Timeout time.Duration `env:"DB_TIMEOUT" envDefault:"60s"`

  dsn     string
  dsnOnce sync.Once
}

func (c *Database) DSN() string {
  c.dsnOnce.Do(func() {
    password := ""
    if len(c.Password) != 0 {
      password = ":" + c.Password
    }
    param := ""
    if len(c.Parameter) != 0 {
      param = "?" + c.Parameter
    }

    c.dsn = fmt.Sprintf("%s://%s%s@%s:%d/%s%s", c.Protocol, c.Username, password,
      c.Host, c.Port, c.Name, param)
  })

  return c.dsn
}
