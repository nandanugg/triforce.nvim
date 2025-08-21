package datapribadi

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

func (s *service) getDataPribadi(ctx context.Context, userID int64) (*dataPribadi, error) {
	return s.repo.getDataPribadi(ctx, userID)
}

func (s *service) listStatusPernikahan(ctx context.Context) ([]statusPernikahan, error) {
	data, err := s.repo.listStatusPernikahan(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	return data, nil
}
