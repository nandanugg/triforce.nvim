package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) getUser(ctx context.Context, id uuid.UUID, source string) (*user, error) {
	query := `select nip from "user" where id = $1 and source = $2 and deleted_at is null`
	var nip string
	if err := r.db.QueryRow(ctx, query, id, source).Scan(&nip); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	query = `
		select distinct on (r.service)
			r.service,
			r.nama
		from user_role ur
		join role r on r.id = ur.role_id and r.deleted_at is null
		where ur.nip = $1 and ur.deleted_at is null
		order by r.service, ur.updated_at desc
	`
	rows, err := r.db.Query(ctx, query, nip)
	if err != nil {
		return nil, fmt.Errorf("query user_role: %w", err)
	}
	defer rows.Close()

	user := &user{nip: nip, roles: make(map[string]string)}
	for rows.Next() {
		var service, nama string
		if err = rows.Scan(&service, &nama); err != nil {
			return nil, fmt.Errorf("row scan user_role: %w", err)
		}

		user.roles[service] = nama
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan user_role: %w", err)
	}

	return user, nil
}
