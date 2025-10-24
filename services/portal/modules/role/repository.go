package role

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type sqlcRepository interface {
	ListRoles(ctx context.Context, arg sqlc.ListRolesParams) ([]sqlc.ListRolesRow, error)
	CountRoles(ctx context.Context) (int64, error)
	GetRole(ctx context.Context, id int16) (sqlc.GetRoleRow, error)
	ListResourcePermissionsByRoleID(ctx context.Context, roleID int16) ([]sqlc.ListResourcePermissionsByRoleIDRow, error)
	ListRoleResourcePermissionsByRoleID(ctx context.Context, roleID int16) ([]sqlc.ListRoleResourcePermissionsByRoleIDRow, error)
	CountResourcePermissionsByIDs(ctx context.Context, ids []int32) (int64, error)

	WithTx(tx pgx.Tx) *sqlc.Queries
}

type txRepository interface {
	CreateRole(ctx context.Context, arg sqlc.CreateRoleParams) (int16, error)
	UpdateRole(ctx context.Context, arg sqlc.UpdateRoleParams) (int16, error)
	CreateRoleResourcePermissions(ctx context.Context, arg sqlc.CreateRoleResourcePermissionsParams) error
	DeleteRoleResourcePermissions(ctx context.Context, arg sqlc.DeleteRoleResourcePermissionsParams) error
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

func (r *repository) create(ctx context.Context, params createParams) (int16, error) {
	var res int16
	err := r.withTransaction(ctx, func(r txRepository) error {
		id, err := r.CreateRole(ctx, sqlc.CreateRoleParams{
			Nama:      params.nama,
			Deskripsi: pgtype.Text{String: params.deskripsi, Valid: true},
			IsDefault: params.isDefault,
		})
		if err != nil {
			return fmt.Errorf("create role: %w", err)
		}

		if err := r.CreateRoleResourcePermissions(ctx, sqlc.CreateRoleResourcePermissionsParams{
			RoleID:                id,
			ResourcePermissionIds: params.resourcePermissionIDs,
		}); err != nil {
			return fmt.Errorf("create role resource permission: %w", err)
		}

		res = id
		return nil
	})
	return res, err
}

func (r *repository) update(ctx context.Context, id int16, opts updateOptions) error {
	return r.withTransaction(ctx, func(r txRepository) error {
		if _, err := r.UpdateRole(ctx, sqlc.UpdateRoleParams{
			ID:        id,
			Nama:      pgtype.Text{String: typeutil.FromPtr(opts.nama), Valid: opts.nama != nil},
			Deskripsi: pgtype.Text{String: typeutil.FromPtr(opts.deskripsi), Valid: opts.deskripsi != nil},
			IsDefault: pgtype.Bool{Bool: typeutil.FromPtr(opts.isDefault), Valid: opts.isDefault != nil},
			IsAktif:   pgtype.Bool{Bool: typeutil.FromPtr(opts.isAktif), Valid: opts.isAktif != nil},
		}); err != nil {
			return fmt.Errorf("update role: %w", err)
		}

		if opts.resourcePermissionIDs != nil {
			if err := r.CreateRoleResourcePermissions(ctx, sqlc.CreateRoleResourcePermissionsParams{
				RoleID:                id,
				ResourcePermissionIds: *opts.resourcePermissionIDs,
			}); err != nil {
				return fmt.Errorf("create role resource permission: %w", err)
			}

			if err := r.DeleteRoleResourcePermissions(ctx, sqlc.DeleteRoleResourcePermissionsParams{
				RoleID:                       id,
				ExcludeResourcePermissionIds: *opts.resourcePermissionIDs,
			}); err != nil {
				return fmt.Errorf("delete role resource permission: %w", err)
			}
		}

		return nil
	})
}
