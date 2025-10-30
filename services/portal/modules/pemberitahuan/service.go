package pemberitahuan

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type repository interface {
	ListPemberitahuan(ctx context.Context, arg sqlc.ListPemberitahuanParams) ([]sqlc.ListPemberitahuanRow, error)
	UpdatePemberitahuan(ctx context.Context, arg sqlc.UpdatePemberitahuanParams) (sqlc.UpdatePemberitahuanRow, error)
	CountPemberitahuan(ctx context.Context, status any) (int64, error)
	CreatePemberitahuan(ctx context.Context, arg sqlc.CreatePemberitahuanParams) (sqlc.CreatePemberitahuanRow, error)
	DeletePemberitahuan(ctx context.Context, id int64) (int64, error)
	GetOverlappingPinnedPemberitahuan(ctx context.Context, arg sqlc.GetOverlappingPinnedPemberitahuanParams) (sqlc.GetOverlappingPinnedPemberitahuanRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, filterPeriode Status, limit, offset uint) ([]pemberitahuan, uint, error) {
	rows, err := s.repo.ListPemberitahuan(ctx, sqlc.ListPemberitahuanParams{
		Limit: int32(limit), Offset: int32(offset),
		Status: filterPeriode,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountPemberitahuan(ctx, filterPeriode)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(r sqlc.ListPemberitahuanRow) pemberitahuan {
		return pemberitahuan{
			ID:                 r.ID,
			JudulBerita:        r.JudulBerita,
			DeskripsiBerita:    r.DeskripsiBerita,
			Pinned:             r.Pinned,
			DiterbitkanPada:    r.DiterbitkanPada.Time,
			DitarikPada:        r.DitarikPada.Time,
			DiperbaruiOleh:     r.UpdatedBy,
			TerakhirDiperbarui: r.UpdatedAt.Time,
			Status:             r.Status,
		}
	}), uint(count), nil
}

type createPemberitahuanParams struct {
	JudulBerita      string
	DeskripsiBerita  string
	Pinned           bool
	DiterbitkanPada  pgtype.Timestamptz
	DitarikPada      pgtype.Timestamptz
	DiperbaharuiOleh string
}

func (s *service) create(ctx context.Context, p createPemberitahuanParams) (*pemberitahuan, error) {
	if p.Pinned {
		if err := s.checkOverlap(ctx, sqlc.GetOverlappingPinnedPemberitahuanParams{
			DitarikPada:     p.DitarikPada,
			DiterbitkanPada: p.DiterbitkanPada,
			ID:              0,
		}); err != nil {
			return nil, err
		}
	}

	r, err := s.repo.CreatePemberitahuan(ctx, sqlc.CreatePemberitahuanParams{
		JudulBerita:     p.JudulBerita,
		DeskripsiBerita: p.DeskripsiBerita,
		Pinned:          p.Pinned,
		DiterbitkanPada: p.DiterbitkanPada,
		DitarikPada:     p.DitarikPada,
		UpdatedBy:       p.DiperbaharuiOleh,
	})
	if err != nil {
		return nil, fmt.Errorf("repo create: %w", err)
	}
	return &pemberitahuan{
		ID:                 int64(r.ID),
		JudulBerita:        r.JudulBerita,
		DeskripsiBerita:    r.DeskripsiBerita,
		Pinned:             r.Pinned,
		DiterbitkanPada:    r.DiterbitkanPada.Time,
		DitarikPada:        r.DitarikPada.Time,
		DiperbaruiOleh:     r.UpdatedBy,
		TerakhirDiperbarui: r.UpdatedAt.Time,
		Status:             r.Status,
	}, nil
}

type updatePemberitahuanParams struct {
	ID               int64
	JudulBerita      string
	DeskripsiBerita  string
	Pinned           bool
	DiterbitkanPada  pgtype.Timestamptz
	DitarikPada      pgtype.Timestamptz
	DiperbaharuiOleh string
}

func (s *service) update(ctx context.Context, id int64, p updatePemberitahuanParams) (*pemberitahuan, error) {
	if p.Pinned {
		if err := s.checkOverlap(ctx, sqlc.GetOverlappingPinnedPemberitahuanParams{
			DitarikPada:     p.DitarikPada,
			DiterbitkanPada: p.DiterbitkanPada,
			ID:              id,
		}); err != nil {
			return nil, err
		}
	}

	r, err := s.repo.UpdatePemberitahuan(ctx, sqlc.UpdatePemberitahuanParams{
		ID:              id,
		JudulBerita:     p.JudulBerita,
		DeskripsiBerita: p.DeskripsiBerita,
		Pinned:          p.Pinned,
		DiterbitkanPada: p.DiterbitkanPada,
		DitarikPada:     p.DitarikPada,
		UpdatedBy:       p.DiperbaharuiOleh,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo update: %w", err)
	}

	return &pemberitahuan{
		ID:                 int64(r.ID),
		JudulBerita:        r.JudulBerita,
		DeskripsiBerita:    r.DeskripsiBerita,
		Pinned:             r.Pinned,
		DiterbitkanPada:    r.DiterbitkanPada.Time,
		DitarikPada:        r.DitarikPada.Time,
		DiperbaruiOleh:     r.UpdatedBy,
		TerakhirDiperbarui: r.UpdatedAt.Time,
		Status:             r.Status,
	}, nil
}

func (s *service) delete(ctx context.Context, id int64) (bool, error) {
	affected, err := s.repo.DeletePemberitahuan(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error: %w", err)
	}
	return affected > 0, nil
}

func (s *service) checkOverlap(ctx context.Context, p sqlc.GetOverlappingPinnedPemberitahuanParams) error {
	overlap, err := s.repo.GetOverlappingPinnedPemberitahuan(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("repo check pinned duplicate: %w", err)
	}
	return NewError(ErrConflict, fmt.Sprintf("conflict with '%s' (id=%d)", overlap.JudulBerita, overlap.ID))
}
