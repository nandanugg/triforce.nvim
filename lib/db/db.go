package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // "pgx" sql driver
)

func New(host, user, password, dbname, schema string) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s search_path=%s",
		host, user, password, dbname, schema,
	))
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := db.PingContext(timeout); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}
