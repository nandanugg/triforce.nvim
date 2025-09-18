package riwayatkinerja

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

func (s *service) list(ctx context.Context, userID int64, limit, offset uint) ([]riwayatKinerja, uint, error) {
	data, err := s.repo.list(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.count(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return data, count, nil
}
