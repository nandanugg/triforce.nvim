package keluarga

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

func (s *service) list(ctx context.Context, userID int64) ([]keluarga, error) {
	ortu, err := s.repo.listOrangTua(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo list orang tua: %w", err)
	}

	pasangan, err := s.repo.listPasangan(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo list pasangan: %w", err)
	}

	anak, err := s.repo.listAnak(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo list anak: %w", err)
	}

	data := ortu
	data = append(data, pasangan...)
	data = append(data, anak...)

	return data, nil
}

func (s *service) listAnak(ctx context.Context, userID int64) ([]keluarga, error) {
	data, err := s.repo.listAnak(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo list anak: %w", err)
	}

	return data, nil
}

func (s *service) listOrangTua(ctx context.Context, userID int64) ([]keluarga, error) {
	data, err := s.repo.listOrangTua(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo list orang tua: %w", err)
	}

	return data, nil
}

func (s *service) listPasangan(ctx context.Context, userID int64) ([]keluarga, error) {
	data, err := s.repo.listPasangan(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo list pasangan: %w", err)
	}

	return data, nil
}
