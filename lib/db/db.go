package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib" // "pgx" sql driver
)

func New(host string, port uint, user, password, dbname, schema string) (*sql.DB, error) {
	connConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		host, port, user, password, dbname,
	))
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, errExec := conn.Exec(ctx, fmt.Sprintf("SET search_path TO %s", schema))
		return errExec
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := db.PingContext(timeout); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}

func NewPgxPool(host string, port uint, user, password, dbname, schema string) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		host, port, user, password, dbname,
	))
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, errExec := conn.Exec(ctx, fmt.Sprintf("SET search_path TO %s", schema))
		return errExec
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return pool, nil
}
