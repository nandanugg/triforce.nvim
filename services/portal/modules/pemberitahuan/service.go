package pemberitahuan

import (
	"context"
	"fmt"
)

type service struct {
	repo *repository
}

func newService(r *repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint, cari string) ([]pemberitahuan, uint, error) {
	data, err := s.repo.list(ctx, limit, offset, cari)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.count(ctx, cari)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return data, count, nil
}
