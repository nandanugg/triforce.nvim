package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type service struct {
	repo *repository
}

func newService(repo *repository) *service {
	return &service{repo: repo}
}

type listOptions struct {
	nip    string
	roleID int16
}

func (s *service) list(ctx context.Context, opts listOptions, limit, offset uint) ([]user, uint, error) {
	rows, err := s.repo.ListUsersGroupByNIP(ctx, sqlc.ListUsersGroupByNIPParams{
		Nip:    pgtype.Text{String: opts.nip, Valid: opts.nip != ""},
		RoleID: pgtype.Int2{Int16: opts.roleID, Valid: opts.roleID != 0},
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list user: %w", err)
	}

	nips := typeutil.Map(rows, func(row sqlc.ListUsersGroupByNIPRow) string {
		return row.Nip
	})

	roles, err := s.repo.ListRolesByNIPs(ctx, nips)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list roles: %w", err)
	}

	rolesMap := typeutil.GroupByMap(roles, func(row sqlc.ListRolesByNIPsRow) (string, role) {
		return row.Nip, role{
			ID:        row.ID,
			Nama:      row.Nama,
			IsDefault: row.IsDefault,
			IsAktif:   row.IsAktif,
		}
	})

	count, err := s.repo.CountUsersGroupByNIP(ctx, sqlc.CountUsersGroupByNIPParams{
		Nip:    pgtype.Text{String: opts.nip, Valid: opts.nip != ""},
		RoleID: pgtype.Int2{Int16: opts.roleID, Valid: opts.roleID != 0},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo count user: %w", err)
	}

	users := make([]user, 0, len(rows))
	for _, row := range rows {
		var profiles []profile
		if err := json.Unmarshal(row.Profiles, &profiles); err != nil {
			return nil, 0, fmt.Errorf("unmarshal profiles: %w", err)
		}

		roles := rolesMap[row.Nip]
		if roles == nil {
			roles = make([]role, 0)
		}

		users = append(users, user{
			NIP:      row.Nip,
			Profiles: profiles,
			Roles:    roles,
		})
	}
	return users, uint(count), nil
}

func (s *service) get(ctx context.Context, nip string) (*user, error) {
	data, err := s.repo.GetUserGroupByNIP(ctx, nip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get user: %w", err)
	}

	roles, err := s.repo.ListRolesByNIPs(ctx, []string{data.Nip})
	if err != nil {
		return nil, fmt.Errorf("repo get roles: %w", err)
	}

	var profiles []profile
	if err := json.Unmarshal(data.Profiles, &profiles); err != nil {
		return nil, fmt.Errorf("unmarshal profiles: %w", err)
	}

	return &user{
		NIP:      data.Nip,
		Profiles: profiles,
		Roles: typeutil.Map(roles, func(row sqlc.ListRolesByNIPsRow) role {
			return role{
				ID:        row.ID,
				Nama:      row.Nama,
				IsDefault: row.IsDefault,
				IsAktif:   row.IsAktif,
			}
		}),
	}, nil
}

func (s *service) update(ctx context.Context, nip string, roleIDs []int16) (bool, error) {
	if err := s.validateRoleIDs(ctx, roleIDs); err != nil {
		return false, err
	}

	found, err := s.repo.IsUserExistsByNIP(ctx, nip)
	if err != nil {
		return false, fmt.Errorf("repo user exists: %w", err)
	}
	if !found {
		return false, nil
	}

	if err := s.repo.update(ctx, nip, roleIDs); err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}
	return true, nil
}

func (s *service) validateRoleIDs(ctx context.Context, ids []int16) error {
	if len(ids) == 0 {
		return nil
	}

	count, err := s.repo.CountRolesByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("repo count roles: %w", err)
	}
	if count != int64(len(ids)) {
		return errRoleNotFound
	}
	return nil
}
