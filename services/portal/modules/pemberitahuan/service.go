package pemberitahuan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
)

type repository interface {
	ListPemberitahuan(ctx context.Context, arg sqlc.ListPemberitahuanParams) ([]sqlc.ListPemberitahuanRow, error)
	ListActivePemberitahuan(ctx context.Context, arg sqlc.ListActivePemberitahuanParams) ([]sqlc.ListActivePemberitahuanRow, error)
	UpdatePemberitahuan(ctx context.Context, arg sqlc.UpdatePemberitahuanParams) (sqlc.UpdatePemberitahuanRow, error)
	CountPemberitahuan(ctx context.Context, arg sqlc.CountPemberitahuanParams) (int64, error)
	CreatePemberitahuan(ctx context.Context, arg sqlc.CreatePemberitahuanParams) (sqlc.CreatePemberitahuanRow, error)
	DeletePemberitahuan(ctx context.Context, id int64) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listActive(ctx context.Context, limit, offset uint) ([]pemberitahuan, uint, error) {
	rows, err := s.repo.ListActivePemberitahuan(ctx, sqlc.ListActivePemberitahuanParams{
		Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountPemberitahuan(ctx, sqlc.CountPemberitahuanParams{
		Status: pemberitahuanStatusActive,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(r sqlc.ListActivePemberitahuanRow) pemberitahuan {
		return pemberitahuan{
			ID:                    r.ID,
			JudulBerita:           r.JudulBerita,
			DeskripsiBerita:       r.DeskripsiBerita,
			PinnedAt:              r.PinnedAt,
			DiterbitkanPada:       r.DiterbitkanPada.Time,
			DitarikPada:           r.DitarikPada.Time,
			DiperbaruiOleh:        r.UpdatedBy,
			TerakhirDiperbarui:    r.UpdatedAt.Time,
			IsCurrentPeriodPinned: &r.IsCurrentPeriodPinned.Bool,
			Status:                "ACTIVE",
		}
	}), uint(count), nil
}

func (s *service) list(ctx context.Context, filterJudul, sortBy string, limit, offset uint) ([]pemberitahuan, uint, error) {
	rows, err := s.repo.ListPemberitahuan(ctx, sqlc.ListPemberitahuanParams{
		Limit: int32(limit), Offset: int32(offset),
		JudulBerita: filterJudul,
		SortBy:      sortBy,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountPemberitahuan(ctx, sqlc.CountPemberitahuanParams{
		Status:      pemberitahuanStatusAll,
		JudulBerita: filterJudul,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(r sqlc.ListPemberitahuanRow) pemberitahuan {
		return pemberitahuan{
			ID:                 r.ID,
			JudulBerita:        r.JudulBerita,
			DeskripsiBerita:    r.DeskripsiBerita,
			PinnedAt:           r.PinnedAt,
			DiterbitkanPada:    r.DiterbitkanPada.Time,
			DitarikPada:        r.DitarikPada.Time,
			DiperbaruiOleh:     r.UpdatedBy,
			TerakhirDiperbarui: r.UpdatedAt.Time,
			Status:             s.statusFromDates(r.DiterbitkanPada.Time, r.DitarikPada.Time),
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
	pinnedAt := pgtype.Timestamptz{Valid: false}
	if p.Pinned {
		pinnedAt.Valid = true
		pinnedAt.Time = time.Now()
	}
	r, err := s.repo.CreatePemberitahuan(ctx, sqlc.CreatePemberitahuanParams{
		JudulBerita:     p.JudulBerita,
		DeskripsiBerita: p.DeskripsiBerita,
		PinnedAt:        pinnedAt,
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
		PinnedAt:           pinnedAt,
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
	pinnedAt := pgtype.Timestamptz{Valid: false}
	if p.Pinned {
		pinnedAt.Valid = true
		pinnedAt.Time = time.Now()
	}
	r, err := s.repo.UpdatePemberitahuan(ctx, sqlc.UpdatePemberitahuanParams{
		ID:              id,
		JudulBerita:     p.JudulBerita,
		PinnedAt:        pinnedAt,
		DeskripsiBerita: p.DeskripsiBerita,
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
		PinnedAt:           pinnedAt,
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

func (s *service) statusFromDates(start, end time.Time) string {
	now := time.Now()

	switch {
	case now.Before(start):
		return "WAITING"
	case (now.Equal(start) || now.After(start)) && now.Before(end):
		return "ACTIVE"
	case now.Equal(end) || now.After(end):
		return "OVER"
	default:
		return "UNKNOWN"
	}
}
