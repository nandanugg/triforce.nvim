package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type service struct {
	repo *repository
}

func newService(repo *repository) *service {
	return &service{repo: repo}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]role, uint, error) {
	list, err := s.repo.ListRoles(ctx, sqlc.ListRolesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRoles(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(list, func(row sqlc.ListRolesRow) role {
		return role{
			ID:         row.ID,
			Nama:       row.Nama,
			Deskripsi:  row.Deskripsi.String,
			IsDefault:  row.IsDefault,
			IsAktif:    row.IsAktif,
			JumlahUser: row.JumlahUser,
		}
	}), uint(count), nil
}

func (s *service) get(ctx context.Context, id int16) (*role, error) {
	data, err := s.repo.GetRole(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get: %w", err)
	}

	resourcePermissionRows, err := s.repo.ListResourcePermissionsByRoleID(ctx, data.ID)
	if err != nil {
		return nil, fmt.Errorf("repo list resource permissions: %w", err)
	}

	resourcePermissions := typeutil.Map(resourcePermissionRows, func(row sqlc.ListResourcePermissionsByRoleIDRow) resourcePermission {
		return resourcePermission{
			ID:             row.ID,
			Kode:           row.Kode.String,
			NamaResource:   row.NamaResource,
			NamaPermission: row.NamaPermission,
		}
	})

	return &role{
		ID:                  data.ID,
		Nama:                data.Nama,
		Deskripsi:           data.Deskripsi.String,
		IsDefault:           data.IsDefault,
		IsAktif:             data.IsAktif,
		JumlahUser:          data.JumlahUser,
		ResourcePermissions: &resourcePermissions,
	}, nil
}

type createParams struct {
	nama                  string
	deskripsi             string
	isDefault             bool
	resourcePermissionIDs []int32
}

func (s *service) create(ctx context.Context, params createParams) (int16, error) {
	if err := s.validateResourcePermissionIDs(ctx, params.resourcePermissionIDs); err != nil {
		return 0, err
	}

	id, err := s.repo.create(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("repo create: %w", err)
	}
	return id, nil
}

type updateOptions struct {
	nama                  *string
	deskripsi             *string
	isDefault             *bool
	isAktif               *bool
	resourcePermissionIDs *[]int32
}

func (s *service) update(ctx context.Context, id int16, opts updateOptions) (bool, error) {
	if opts.resourcePermissionIDs != nil {
		if err := s.validateResourcePermissionIDs(ctx, *opts.resourcePermissionIDs); err != nil {
			return false, err
		}
	}

	if err := s.repo.update(ctx, id, opts); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repo update: %w", err)
	}
	return true, nil
}

func (s *service) validateResourcePermissionIDs(ctx context.Context, ids []int32) error {
	if len(ids) == 0 {
		return nil
	}

	count, err := s.repo.CountResourcePermissionsByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("repo count resource permissions: %w", err)
	}
	if count != int64(len(ids)) {
		return errResourcePermissionNotFound
	}
	return nil
}
