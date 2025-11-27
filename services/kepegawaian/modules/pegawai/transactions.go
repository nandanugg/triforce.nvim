package pegawai

import (
	"context"
)

func (s *service) withTransaction(
	ctx context.Context,
	fn func(txRepo repository) error,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txRepo := s.repo.WithTx(tx)

	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
