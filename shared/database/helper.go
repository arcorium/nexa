package database

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"nexa/shared/variadic"
)

func OpenPostgres(config *Config, log bool, models ...any) (*bun.DB, error) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithTimeout(config.Timeout)))
	// Test connection
	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqlDb, pgdialect.New())
	if log {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	db.RegisterModel(variadic.New(models...).Values())

	return db, nil
}

func Seed[T any](db bun.IDB, values ...T) error {
	ctx := context.Background()
	_, err := db.NewInsert().
		Model(&values).
		Exec(ctx)

	return err
}
