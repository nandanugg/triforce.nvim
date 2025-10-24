package user

import (
	"context"
	"fmt"
	"log/slog"

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

type txRepository interface {
	CreateUserRoles(ctx context.Context, arg sqlc.CreateUserRolesParams) error
	DeleteUserRoles(ctx context.Context, arg sqlc.DeleteUserRolesParams) error
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

func (r *repository) withTransaction(ctx context.Context, fn func(txRepository) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			slog.WarnContext(ctx, "Error rollback transaction", "error", err)
		}
	}()

	if err := fn(r.WithTx(tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *repository) update(ctx context.Context, nip string, roleIDs []int16) error {
	return r.withTransaction(ctx, func(r txRepository) error {
		if err := r.CreateUserRoles(ctx, sqlc.CreateUserRolesParams{
			Nip:     nip,
			RoleIds: roleIDs,
		}); err != nil {
			return fmt.Errorf("create user role: %w", err)
		}

		if err := r.DeleteUserRoles(ctx, sqlc.DeleteUserRolesParams{
			Nip:            nip,
			ExcludeRoleIds: roleIDs,
		}); err != nil {
			return fmt.Errorf("delete user role: %w", err)
		}

		return nil
	})
}
