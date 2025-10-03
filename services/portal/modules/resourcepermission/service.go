package resourcepermission

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type repository interface {
	ListResources(ctx context.Context, arg sqlc.ListResourcesParams) ([]sqlc.ListResourcesRow, error)
	CountResources(ctx context.Context) (int64, error)
	ListResourcePermissionsByResourceIDs(ctx context.Context, resourceids []int16) ([]sqlc.ListResourcePermissionsByResourceIDsRow, error)
	ListResourcePermissionsByNip(ctx context.Context, nip string) ([]pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) listResourcePermissionsByNip(ctx context.Context, nip string) ([]string, error) {
	list, err := s.repo.ListResourcePermissionsByNip(ctx, nip)
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	return typeutil.Map(list, func(item pgtype.Text) string {
		return item.String
	}), nil
}

func (s *service) listResources(ctx context.Context, limit, offset uint) ([]resource, uint, error) {
	resources, err := s.repo.ListResources(ctx, sqlc.ListResourcesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list resources: %w", err)
	}

	resourceIDs := typeutil.Map(resources, func(row sqlc.ListResourcesRow) int16 {
		return row.ID
	})
	resourcePermissions, err := s.repo.ListResourcePermissionsByResourceIDs(ctx, resourceIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list resource permissions: %w", err)
	}

	resourcePermissionsMap := make(map[int16][]resourcePermission, len(resources))
	for _, row := range resourcePermissions {
		resourcePermissionsMap[row.ResourceID] = append(resourcePermissionsMap[row.ResourceID], resourcePermission{
			ID:             row.ID,
			Kode:           row.Kode.String,
			NamaPermission: row.NamaPermission,
		})
	}

	total, err := s.repo.CountResources(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count resources: %w", err)
	}

	return typeutil.Map(resources, func(row sqlc.ListResourcesRow) resource {
		return resource{
			Nama:                row.Nama,
			ResourcePermissions: resourcePermissionsMap[row.ID],
		}
	}), uint(total), nil
}
