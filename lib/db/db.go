package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	PingTimeout time.Duration
}

func New(host string, port uint, user, password, dbname, schema string, opts ...Options) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		host, port, user, password, dbname,
	))
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, errExec := conn.Exec(ctx, "SET search_path TO "+schema)
		return errExec
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	opt := Options{
		PingTimeout: time.Second,
	}
	if len(opts) > 0 && opts[0].PingTimeout > 0 {
		opt.PingTimeout = opts[0].PingTimeout
	}
	timeout, cancel := context.WithTimeout(context.Background(), opt.PingTimeout)
	defer cancel()
	if err := pool.Ping(timeout); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return pool, nil
}
