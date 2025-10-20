package template

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repo interface {
	CountTemplates(ctx context.Context) (int64, error)
	ListTemplates(ctx context.Context, arg sqlc.ListTemplatesParams) ([]sqlc.ListTemplatesRow, error)
	GetTemplate(ctx context.Context, id int32) (sqlc.GetTemplateRow, error)
	GetTemplateBerkas(ctx context.Context, id int32) (pgtype.Text, error)
	CreateTemplate(ctx context.Context, arg sqlc.CreateTemplateParams) (sqlc.CreateTemplateRow, error)
	UpdateTemplate(ctx context.Context, arg sqlc.UpdateTemplateParams) (sqlc.UpdateTemplateRow, error)
	DeleteTemplate(ctx context.Context, id int32) (int64, error)
}

type service struct {
	repo repo
}

func newService(r repo) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]template, int64, error) {
	rows, err := s.repo.ListTemplates(ctx, sqlc.ListTemplatesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error ListTemplate: %w", err)
	}

	total, err := s.repo.CountTemplates(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error CountTemplate: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListTemplatesRow) template {
		return template{
			ID:        row.ID,
			Nama:      row.Nama,
			Filename:  row.Filename,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}), total, nil
}

func (s *service) get(ctx context.Context, id int32) (*template, error) {
	row, err := s.repo.GetTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetTemplate: %w", err)
	}

	result := &template{
		ID:        row.ID,
		Nama:      row.Nama,
		Filename:  row.Filename,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	return result, nil
}

func (s *service) getBerkas(ctx context.Context, id int32) (string, []byte, error) {
	row, err := s.repo.GetTemplateBerkas(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, nil
		}
		return "", nil, fmt.Errorf("[get] error GetTemplateBerkas: %w", err)
	}

	return api.GetMimeTypeAndDecodedData(row.String)
}

type createParams struct {
	nama     string
	filename string
	file     string
}

func (s *service) create(ctx context.Context, params createParams) (*template, error) {
	row, err := s.repo.CreateTemplate(ctx, sqlc.CreateTemplateParams{
		Name:       params.nama,
		Filename:   params.filename,
		FileBase64: params.file,
	})
	if err != nil {
		return nil, fmt.Errorf("[create] error CreateTemplate: %w", err)
	}

	result := &template{
		ID:        row.ID,
		Nama:      row.Nama,
		Filename:  row.Filename,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return result, nil
}

type updateParams struct {
	nama     string
	filename string
	file     string
}

func (s *service) update(ctx context.Context, id int32, params updateParams) (*template, error) {
	row, err := s.repo.UpdateTemplate(ctx, sqlc.UpdateTemplateParams{
		ID:         id,
		Name:       params.nama,
		Filename:   params.filename,
		FileBase64: params.file,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error UpdateTemplate: %w", err)
	}

	result := &template{
		ID:        row.ID,
		Nama:      row.Nama,
		Filename:  row.Filename,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteTemplate(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error DeleteTemplate: %w", err)
	}
	return affected > 0, nil
}
