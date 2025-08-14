package datapribadi

import "context"

type service struct {
	repo *repository
}

func newService(r *repository) *service {
	return &service{repo: r}
}

func (s *service) getDataPribadi(ctx context.Context, userID int64) (*dataPribadi, error) {
	return s.repo.getDataPribadi(ctx, userID)
}
