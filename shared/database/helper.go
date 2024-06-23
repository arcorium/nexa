package database

import (
  "context"
  "database/sql"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/dialect/pgdialect"
  "github.com/uptrace/bun/driver/pgdriver"
  "github.com/uptrace/bun/extra/bundebug"
  "nexa/shared/config"
)

func OpenPostgres(config *config.Database, log bool) (*bun.DB, error) {
  options := []pgdriver.Option{
    pgdriver.WithDSN(config.DSN()),
    pgdriver.WithInsecure(!config.IsSecure),
  }

  if config.Timeout.Milliseconds() > 0 {
    options = append(options, pgdriver.WithTimeout(config.Timeout))
  }

  sqlDb := sql.OpenDB(pgdriver.NewConnector(options...))
  // Test connection
  if err := sqlDb.Ping(); err != nil {
    return nil, err
  }

  db := bun.NewDB(sqlDb, pgdialect.New())
  if log {
    db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
  }
  return db, nil
}

func Seed[T any](db bun.IDB, values ...T) error {
  ctx := context.Background()
  _, err := db.NewInsert().
    Model(&values).
    Exec(ctx)

  return err
}
