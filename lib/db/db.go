package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	PingTimeout time.Duration
}

func New(host string, port uint, user, password, dbname, schema string, opts ...Options) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s options='-c search_path=%s'",
		host, port, user, password, dbname, schema,
	))
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingTimeout := time.Second
	if len(opts) > 0 && opts[0].PingTimeout > 0 {
		pingTimeout = opts[0].PingTimeout
	}
	timeout, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	if err := pool.Ping(timeout); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return pool, nil
}
