package role

import (
	"context"
	"fmt"

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

func (r *repository) create(ctx context.Context, params createParams) (int16, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	qtx := r.WithTx(tx)
	id, err := qtx.CreateRole(ctx, sqlc.CreateRoleParams{
		Nama:      params.nama,
		Deskripsi: pgtype.Text{String: params.deskripsi, Valid: true},
		IsDefault: params.isDefault,
	})
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, fmt.Errorf("create role: %w", err)
	}

	if err := qtx.CreateRoleResourcePermissions(ctx, sqlc.CreateRoleResourcePermissionsParams{
		RoleID:                id,
		ResourcePermissionIds: params.resourcePermissionIDs,
	}); err != nil {
		_ = tx.Rollback(ctx)
		return 0, fmt.Errorf("create role resource permission: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return id, nil
}

func (r *repository) update(ctx context.Context, id int16, opts updateOptions) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	qtx := r.WithTx(tx)
	if _, err := qtx.UpdateRole(ctx, sqlc.UpdateRoleParams{
		ID:        id,
		Nama:      pgtype.Text{String: typeutil.FromPtr(opts.nama), Valid: opts.nama != nil},
		Deskripsi: pgtype.Text{String: typeutil.FromPtr(opts.deskripsi), Valid: opts.deskripsi != nil},
		IsDefault: pgtype.Bool{Bool: typeutil.FromPtr(opts.isDefault), Valid: opts.isDefault != nil},
		IsAktif:   pgtype.Bool{Bool: typeutil.FromPtr(opts.isAktif), Valid: opts.isAktif != nil},
	}); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("update role: %w", err)
	}

	if opts.resourcePermissionIDs != nil {
		if err := qtx.CreateRoleResourcePermissions(ctx, sqlc.CreateRoleResourcePermissionsParams{
			RoleID:                id,
			ResourcePermissionIds: *opts.resourcePermissionIDs,
		}); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("create role resource permission: %w", err)
		}

		if err := qtx.DeleteRoleResourcePermissions(ctx, sqlc.DeleteRoleResourcePermissionsParams{
			RoleID:                       id,
			ExcludeResourcePermissionIds: *opts.resourcePermissionIDs,
		}); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("delete role resource permission: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
