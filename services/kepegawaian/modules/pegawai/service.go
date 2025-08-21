package pegawai

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

type listOptions struct {
	cari       string
	unitID     string
	golonganID int64
	jabatanID  string
	status     string
}

func (s *service) list(ctx context.Context, limit, offset uint64, opts listOptions) ([]pegawai, uint, error) {
	data, err := s.repo.list(ctx, limit, offset, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.count(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return data, count, nil
}

func (s *service) listStatusPegawai(ctx context.Context) ([]statusPegawai, error) {
	data, err := s.repo.listStatusPegawai(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	return data, nil
}
