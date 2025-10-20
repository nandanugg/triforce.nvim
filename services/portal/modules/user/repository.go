package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type sqlcRepository interface {
	GetUserGroupByNIP(ctx context.Context, nip string) (sqlc.GetUserGroupByNIPRow, error)
	ListUsersGroupByNIP(ctx context.Context, arg sqlc.ListUsersGroupByNIPParams) ([]sqlc.ListUsersGroupByNIPRow, error)
	CountUsersGroupByNIP(ctx context.Context, arg sqlc.CountUsersGroupByNIPParams) (int64, error)
	ListRolesByNIPs(ctx context.Context, nips []string) ([]sqlc.ListRolesByNIPsRow, error)
	IsUserExistsByNIP(ctx context.Context, nip string) (bool, error)
	CountRolesByIDs(ctx context.Context, ids []int16) (int64, error)

	WithTx(tx pgx.Tx) *sqlc.Queries
}

type repository struct {
	db *pgxpool.Pool
	sqlcRepository
}

func newRepository(db *pgxpool.Pool, repo sqlcRepository) *repository {
	return &repository{
		db:             db,
		sqlcRepository: repo,
	}
}

func (r *repository) update(ctx context.Context, nip string, roleIDs []int16) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	qtx := r.WithTx(tx)
	if err := qtx.CreateUserRoles(ctx, sqlc.CreateUserRolesParams{
		Nip:     nip,
		RoleIds: roleIDs,
	}); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("create user role: %w", err)
	}

	if err := qtx.DeleteUserRoles(ctx, sqlc.DeleteUserRolesParams{
		Nip:            nip,
		ExcludeRoleIds: roleIDs,
	}); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("delete user role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
