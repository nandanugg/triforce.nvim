package riwayatpenghargaan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	upd "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
)

type repository interface {
	ListRiwayatPenghargaan(ctx context.Context, arg repo.ListRiwayatPenghargaanParams) ([]repo.ListRiwayatPenghargaanRow, error)
	CountRiwayatPenghargaan(ctx context.Context, nip string) (int64, error)
	GetBerkasRiwayatPenghargaan(ctx context.Context, arg repo.GetBerkasRiwayatPenghargaanParams) (pgtype.Text, error)
	GetRiwayatPenghargaan(ctx context.Context, arg repo.GetRiwayatPenghargaanParams) (repo.GetRiwayatPenghargaanRow, error)
	DeleteRiwayatPenghargaan(ctx context.Context, arg repo.DeleteRiwayatPenghargaanParams) (int64, error)
	CreateRiwayatPenghargaan(ctx context.Context, arg repo.CreateRiwayatPenghargaanParams) (int32, error)
	UpdateRiwayatPenghargaan(ctx context.Context, arg repo.UpdateRiwayatPenghargaanParams) (int64, error)
	UpdateRiwayatPenghargaanBerkas(ctx context.Context, arg repo.UpdateRiwayatPenghargaanBerkasParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPenghargaan, uint, error) {
	data, err := s.repo.ListRiwayatPenghargaan(ctx, repo.ListRiwayatPenghargaanParams{
		Nip:    nip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPenghargaan(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatPenghargaanRow) riwayatPenghargaan {
		return riwayatPenghargaan{
			ID:               int(row.ID),
			JenisPenghargaan: row.JenisPenghargaan.String,
			NamaPenghargaan:  row.NamaPenghargaan.String,
			Deskripsi:        row.DeskripsiPenghargaan.String,
			Tanggal:          db.Date(row.TanggalPenghargaan.Time),
		}
	})

	return result, uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPenghargaan(ctx, repo.GetBerkasRiwayatPenghargaanParams{
		Nip: pgNip,
		ID:  id,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}
	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}

func (s *service) uploadBerkas(ctx context.Context, id int32, nip string, base64 string) (bool, error) {
	res, err := s.repo.UpdateRiwayatPenghargaanBerkas(ctx, repo.UpdateRiwayatPenghargaanBerkasParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{Valid: true, String: base64},
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("repo upload berkas: %w", err)
	}

	if res == 0 {
		return false, nil
	}

	return true, nil
}

func (s *service) create(ctx context.Context, nip string, params upsertParams) (int32, error) {
	_, valid := s.validateJenisPenghargaan(params.JenisPenghargaan)
	if !valid {
		return 0, NewError(ErrJenisPenghargaanInvalid, params.JenisPenghargaan)
	}

	id, err := s.repo.CreateRiwayatPenghargaan(ctx, repo.CreateRiwayatPenghargaanParams{
		Nip:                  pgtype.Text{Valid: true, String: nip},
		NamaPenghargaan:      pgtype.Text{Valid: true, String: params.NamaPenghargaan},
		JenisPenghargaan:     pgtype.Text{Valid: true, String: params.JenisPenghargaan},
		DeskripsiPenghargaan: pgtype.Text{Valid: params.Deskripsi != "", String: params.Deskripsi},
		TanggalPenghargaan:   pgtype.Date{Valid: params.Tanggal.ToPgtypeDate().Valid, Time: params.Tanggal.ToPgtypeDate().Time},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id int32, nip string, params upsertParams) (bool, error) {
	_, valid := s.validateJenisPenghargaan(params.JenisPenghargaan)
	if !valid {
		return false, NewError(ErrJenisPenghargaanInvalid, params.JenisPenghargaan)
	}

	res, err := s.repo.UpdateRiwayatPenghargaan(ctx, repo.UpdateRiwayatPenghargaanParams{
		ID:                   id,
		Nip:                  pgtype.Text{Valid: true, String: nip},
		NamaPenghargaan:      pgtype.Text{Valid: true, String: params.NamaPenghargaan},
		JenisPenghargaan:     pgtype.Text{Valid: true, String: params.JenisPenghargaan},
		DeskripsiPenghargaan: pgtype.Text{Valid: params.Deskripsi != "", String: params.Deskripsi},
		TanggalPenghargaan:   pgtype.Date{Valid: params.Tanggal.ToPgtypeDate().Valid, Time: params.Tanggal.ToPgtypeDate().Time},
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	if res == 0 {
		return false, nil
	}

	return true, nil
}

func (s *service) delete(ctx context.Context, id int32, nip string) (bool, error) {
	res, err := s.repo.DeleteRiwayatPenghargaan(ctx, repo.DeleteRiwayatPenghargaanParams{
		ID: id, Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	if res == 0 {
		return false, nil
	}
	return true, nil
}

// GeneratePerubahanData implements usulanperubahandata.ServiceInterface
func (s *service) GeneratePerubahanData(ctx context.Context, nip, action, dataID string, requestData json.RawMessage) ([]byte, error) {
	var data usulanPerubahanData
	if action == upd.ActionUpdate || action == upd.ActionDelete {
		id, err := strconv.ParseInt(dataID, 10, 32)
		if err != nil {
			return nil, api.NewMultiError([]error{errors.New("invalid data ID")})
		}

		prevData, err := s.repo.GetRiwayatPenghargaan(ctx, repo.GetRiwayatPenghargaanParams{
			Nip: nip,
			ID:  int32(id),
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, api.NewMultiError([]error{errors.New("data riwayat penghargaan tidak ditemukan")})
			}
			return nil, err
		}

		data.JenisPenghargaan[0] = prevData.JenisPenghargaan
		data.NamaPenghargaan[0] = prevData.NamaPenghargaan
		data.Deskripsi[0] = prevData.DeskripsiPenghargaan
		data.Tanggal[0] = db.Date(prevData.TanggalPenghargaan.Time)
	}

	if action == upd.ActionCreate || action == upd.ActionUpdate {
		var req upsertParams
		if err := json.Unmarshal(requestData, &req); err != nil {
			return nil, err
		}

		_, valid := s.validateJenisPenghargaan(req.JenisPenghargaan)
		if !valid {
			return nil, NewError(ErrJenisPenghargaanInvalid, req.JenisPenghargaan)
		}

		data.JenisPenghargaan[1] = pgtype.Text{String: req.JenisPenghargaan, Valid: true}
		data.NamaPenghargaan[1] = pgtype.Text{String: req.NamaPenghargaan, Valid: true}
		data.Deskripsi[1] = pgtype.Text{String: req.Deskripsi, Valid: req.Deskripsi != ""}
		data.Tanggal[1] = req.Tanggal
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SyncPerubahanData implements usulanperubahandata.ServiceInterface
func (*service) SyncPerubahanData(ctx context.Context, sqlcTx *repo.Queries, nip, action, dataID string, perubahanData []byte) error {
	var data usulanPerubahanData
	if action == upd.ActionCreate || action == upd.ActionUpdate {
		if err := json.Unmarshal(perubahanData, &data); err != nil {
			return fmt.Errorf("json unmarshal: %w", err)
		}
	}

	switch action {
	case upd.ActionCreate:
		if _, err := sqlcTx.CreateRiwayatPenghargaan(ctx, repo.CreateRiwayatPenghargaanParams{
			Nip:                  pgtype.Text{String: nip, Valid: true},
			NamaPenghargaan:      data.NamaPenghargaan[1],
			JenisPenghargaan:     data.JenisPenghargaan[1],
			DeskripsiPenghargaan: data.Deskripsi[1],
			TanggalPenghargaan:   data.Tanggal[1].ToPgtypeDate(),
		}); err != nil {
			return fmt.Errorf("repo create: %w", err)
		}

	case upd.ActionUpdate:
		id, err := strconv.ParseInt(dataID, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid data ID: %w", err)
		}

		if _, err := sqlcTx.UpdateRiwayatPenghargaan(ctx, repo.UpdateRiwayatPenghargaanParams{
			ID:                   int32(id),
			Nip:                  pgtype.Text{String: nip, Valid: true},
			NamaPenghargaan:      data.NamaPenghargaan[1],
			JenisPenghargaan:     data.JenisPenghargaan[1],
			DeskripsiPenghargaan: data.Deskripsi[1],
			TanggalPenghargaan:   data.Tanggal[1].ToPgtypeDate(),
		}); err != nil {
			return fmt.Errorf("repo update: %w", err)
		}

	case upd.ActionDelete:
		id, err := strconv.ParseInt(dataID, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid data ID: %w", err)
		}

		if _, err := sqlcTx.DeleteRiwayatPenghargaan(ctx, repo.DeleteRiwayatPenghargaanParams{
			ID:  int32(id),
			Nip: nip,
		}); err != nil {
			return fmt.Errorf("repo delete: %w", err)
		}

	default:
		return errors.New("unimplemented action")
	}

	return nil
}
