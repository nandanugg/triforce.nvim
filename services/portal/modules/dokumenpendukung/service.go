package dokumenpendukung

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

func (s *service) list(ctx context.Context) ([]dokumenPendukung, error) {
	data, err := s.repo.list(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	return data, nil
}
