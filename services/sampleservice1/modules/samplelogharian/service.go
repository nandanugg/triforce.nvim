package samplelogharian

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

func (s *service) list(ctx context.Context, userID string, limit, offset uint) ([]logHarian, uint, error) {
	lhs, err := s.repo.list(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.listCount(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list count: %w", err)
	}

	return lhs, count, nil
}
