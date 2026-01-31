package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// TODO: add to golib

func Connect(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("empty db dsn")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err = db.PingContext(pingCtx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB, dir string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}

	if dir == "" {
		return fmt.Errorf("nil dir")
	}

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	migCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return goose.UpContext(migCtx, db, dir)
}
