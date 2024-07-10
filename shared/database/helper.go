package database

import (
  "context"
  "database/sql"
  "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/types"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/dialect/pgdialect"
  "github.com/uptrace/bun/driver/pgdriver"
  "github.com/uptrace/bun/extra/bundebug"
  "time"
)

func OpenPostgresWithConfig(config *config.PostgresDatabase, log bool) (*bun.DB, error) {
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

func OpenPostgres(dsn string, secure bool, timeout time.Duration, log bool) (*bun.DB, error) {
  options := []pgdriver.Option{
    pgdriver.WithDSN(dsn),
    pgdriver.WithInsecure(!secure),
  }

  if timeout.Milliseconds() > 0 {
    options = append(options, pgdriver.WithTimeout(timeout))
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
  res, err := db.NewInsert().
    Model(&values).
    Returning("NULL").
    Exec(ctx)

  if err != nil {
    return err
  }

  if types.Must(res.RowsAffected()) == 0 {
    return sql.ErrNoRows
  }

  return nil
}
