package riwayatpendidikan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	upd "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
)

type repository interface {
	CountRiwayatPendidikan(ctx context.Context, nip pgtype.Text) (int64, error)
	ListRiwayatPendidikan(ctx context.Context, arg sqlc.ListRiwayatPendidikanParams) ([]sqlc.ListRiwayatPendidikanRow, error)
	GetBerkasRiwayatPendidikan(ctx context.Context, arg sqlc.GetBerkasRiwayatPendidikanParams) (pgtype.Text, error)
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	GetRefPendidikan(ctx context.Context, id string) (sqlc.GetRefPendidikanRow, error)
	GetRefTingkatPendidikan(ctx context.Context, id int32) (sqlc.GetRefTingkatPendidikanRow, error)
	GetRiwayatPendidikan(ctx context.Context, arg sqlc.GetRiwayatPendidikanParams) (sqlc.GetRiwayatPendidikanRow, error)

	CreateRiwayatPendidikan(ctx context.Context, arg sqlc.CreateRiwayatPendidikanParams) (int32, error)
	UpdateRiwayatPendidikan(ctx context.Context, arg sqlc.UpdateRiwayatPendidikanParams) (int64, error)
	UploadBerkasRiwayatPendidikan(ctx context.Context, arg sqlc.UploadBerkasRiwayatPendidikanParams) (int64, error)
	DeleteRiwayatPendidikan(ctx context.Context, arg sqlc.DeleteRiwayatPendidikanParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPendidikan, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatPendidikan(ctx, sqlc.ListRiwayatPendidikanParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPendidikan(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatPendidikanRow) riwayatPendidikan {
		return riwayatPendidikan{
			ID:                   row.ID,
			TingkatPendidikanID:  row.TingkatPendidikanID,
			JenjangPendidikan:    row.JenjangPendidikan.String,
			PendidikanID:         row.PendidikanID,
			Pendidikan:           row.Pendidikan.String,
			NamaSekolah:          row.NamaSekolah.String,
			TahunLulus:           row.TahunLulus,
			NomorIjazah:          row.NoIjazah.String,
			GelarDepan:           row.GelarDepan.String,
			GelarBelakang:        row.GelarBelakang.String,
			TugasBelajar:         labelStatusBelajar[row.TugasBelajar.Int16],
			KeteranganPendidikan: row.NegaraSekolah.String,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPendidikan(ctx, sqlc.GetBerkasRiwayatPendidikanParams{
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

func (s *service) create(ctx context.Context, nip string, params upsertParams) (int32, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	if _, err := s.validateReferences(ctx, params); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateRiwayatPendidikan(ctx, sqlc.CreateRiwayatPendidikanParams{
		TingkatPendidikanID: pgtype.Int2{Int16: params.TingkatPendidikanID, Valid: true},
		PendidikanID:        pgtype.Text{String: typeutil.FromPtr(params.PendidikanID), Valid: params.PendidikanID != nil},
		NamaSekolah:         pgtype.Text{String: params.NamaSekolah, Valid: true},
		TahunLulus:          pgtype.Int2{Int16: params.TahunLulus, Valid: true},
		NoIjazah:            pgtype.Text{String: params.NomorIjazah, Valid: true},
		GelarDepan:          pgtype.Text{String: params.GelarDepan, Valid: params.GelarDepan != ""},
		GelarBelakang:       pgtype.Text{String: params.GelarBelakang, Valid: params.GelarBelakang != ""},
		NegaraSekolah:       pgtype.Text{String: params.NegaraSekolah, Valid: params.NegaraSekolah != ""},
		TugasBelajar:        params.TugasBelajar.toID(),
		PnsID:               pgtype.Text{String: pnsID, Valid: true},
		Nip:                 pgtype.Text{String: nip, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id int32, nip string, params upsertParams) (bool, error) {
	if _, err := s.validateReferences(ctx, params); err != nil {
		return false, err
	}

	affected, err := s.repo.UpdateRiwayatPendidikan(ctx, sqlc.UpdateRiwayatPendidikanParams{
		ID:                  id,
		Nip:                 nip,
		TingkatPendidikanID: pgtype.Int2{Int16: params.TingkatPendidikanID, Valid: true},
		PendidikanID:        pgtype.Text{String: typeutil.FromPtr(params.PendidikanID), Valid: params.PendidikanID != nil},
		NamaSekolah:         pgtype.Text{String: params.NamaSekolah, Valid: true},
		TahunLulus:          pgtype.Int2{Int16: params.TahunLulus, Valid: true},
		NoIjazah:            pgtype.Text{String: params.NomorIjazah, Valid: true},
		GelarDepan:          pgtype.Text{String: params.GelarDepan, Valid: params.GelarDepan != ""},
		GelarBelakang:       pgtype.Text{String: params.GelarBelakang, Valid: params.GelarBelakang != ""},
		NegaraSekolah:       pgtype.Text{String: params.NegaraSekolah, Valid: params.NegaraSekolah != ""},
		TugasBelajar:        params.TugasBelajar.toID(),
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id int32, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatPendidikan(ctx, sqlc.DeleteRiwayatPendidikanParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int32, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatPendidikan(ctx, sqlc.UploadBerkasRiwayatPendidikanParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo upload berkas: %w", err)
	}

	return affected > 0, nil
}

type references struct {
	tingkatPendidikan sqlc.GetRefTingkatPendidikanRow
	pendidikan        sqlc.GetRefPendidikanRow
}

func (s *service) validateReferences(ctx context.Context, params upsertParams) (*references, error) {
	var errs []error
	var err error

	var ref references
	if ref.tingkatPendidikan, err = s.repo.GetRefTingkatPendidikan(ctx, int32(params.TingkatPendidikanID)); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get tingkat pendidikan: %w", err)
		}
		errs = append(errs, errTingkatPendidikanNotFound)
	}

	if params.PendidikanID != nil {
		if ref.pendidikan, err = s.repo.GetRefPendidikan(ctx, *params.PendidikanID); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("repo get pendidikan: %w", err)
			}
			errs = append(errs, errPendidikanNotFound)
		}
	}

	if len(errs) > 0 {
		return nil, api.NewMultiError(errs)
	}
	return &ref, nil
}

// GeneratePerubahanData implements usulanperubahandata.ServiceInterface
func (s *service) GeneratePerubahanData(ctx context.Context, nip, action, dataID string, requestData json.RawMessage) ([]byte, error) {
	var data usulanPerubahanData
	if action == upd.ActionUpdate || action == upd.ActionDelete {
		id, err := strconv.ParseInt(dataID, 10, 32)
		if err != nil {
			return nil, api.NewMultiError([]error{errors.New(`parameter "data_id" harus dalam format yang sesuai`)})
		}

		prevData, err := s.repo.GetRiwayatPendidikan(ctx, sqlc.GetRiwayatPendidikanParams{
			Nip: nip,
			ID:  int32(id),
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, api.NewMultiError([]error{errors.New("data riwayat pendidikan tidak ditemukan")})
			}
			return nil, err
		}

		data.TingkatPendidikanID[0] = prevData.TingkatPendidikanID
		data.TingkatPendidikan[0] = prevData.TingkatPendidikan
		data.PendidikanID[0] = prevData.PendidikanID
		data.Pendidikan[0] = prevData.Pendidikan
		data.NamaSekolah[0] = prevData.NamaSekolah
		data.TahunLulus[0] = prevData.TahunLulus
		data.NomorIjazah[0] = prevData.NoIjazah
		data.GelarDepan[0] = prevData.GelarDepan
		data.GelarBelakang[0] = prevData.GelarBelakang
		data.TugasBelajar[0] = pgtype.Text{
			String: string(labelStatusBelajar[prevData.TugasBelajar.Int16]),
			Valid:  prevData.TugasBelajar.Valid,
		}
		data.NegaraSekolah[0] = prevData.NegaraSekolah
	}

	if action == upd.ActionCreate || action == upd.ActionUpdate {
		var req upsertParams
		if err := json.Unmarshal(requestData, &req); err != nil {
			return nil, err
		}

		ref, err := s.validateReferences(ctx, req)
		if err != nil {
			return nil, err
		}

		data.TingkatPendidikanID[1] = pgtype.Int2{Int16: req.TingkatPendidikanID, Valid: true}
		data.TingkatPendidikan[1] = ref.tingkatPendidikan.Nama
		data.PendidikanID[1] = pgtype.Text{String: typeutil.FromPtr(req.PendidikanID), Valid: req.PendidikanID != nil}
		data.Pendidikan[1] = ref.pendidikan.Nama
		data.NamaSekolah[1] = pgtype.Text{String: req.NamaSekolah, Valid: true}
		data.TahunLulus[1] = pgtype.Int2{Int16: req.TahunLulus, Valid: true}
		data.NomorIjazah[1] = pgtype.Text{String: req.NomorIjazah, Valid: true}
		data.GelarDepan[1] = pgtype.Text{String: req.GelarDepan, Valid: req.GelarDepan != ""}
		data.GelarBelakang[1] = pgtype.Text{String: req.GelarBelakang, Valid: req.GelarBelakang != ""}
		data.TugasBelajar[1] = pgtype.Text{String: string(req.TugasBelajar), Valid: req.TugasBelajar.toID().Valid}
		data.NegaraSekolah[1] = pgtype.Text{String: req.NegaraSekolah, Valid: req.NegaraSekolah != ""}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SyncPerubahanData implements usulanperubahandata.ServiceInterface
func (*service) SyncPerubahanData(ctx context.Context, sqlcTx *sqlc.Queries, nip, action, dataID string, perubahanData []byte) error {
	var id int32
	if action == upd.ActionUpdate || action == upd.ActionDelete {
		idInt, err := strconv.ParseInt(dataID, 10, 32)
		if err != nil {
			return fmt.Errorf("parse id: %w", err)
		}
		id = int32(idInt)
	}

	var data usulanPerubahanData
	if action == upd.ActionCreate || action == upd.ActionUpdate {
		if err := json.Unmarshal(perubahanData, &data); err != nil {
			return fmt.Errorf("json unmarshal: %w", err)
		}
	}

	switch action {
	case upd.ActionCreate:
		pnsID, err := sqlcTx.GetPegawaiPNSIDByNIP(ctx, nip)
		if err != nil {
			return fmt.Errorf("repo get pegawai: %w", err)
		}

		if _, err := sqlcTx.CreateRiwayatPendidikan(ctx, sqlc.CreateRiwayatPendidikanParams{
			TingkatPendidikanID: data.TingkatPendidikanID[1],
			PendidikanID:        data.PendidikanID[1],
			NamaSekolah:         data.NamaSekolah[1],
			TahunLulus:          data.TahunLulus[1],
			NoIjazah:            data.NomorIjazah[1],
			GelarDepan:          data.GelarDepan[1],
			GelarBelakang:       data.GelarBelakang[1],
			NegaraSekolah:       data.NegaraSekolah[1],
			TugasBelajar:        statusBelajar(data.TugasBelajar[1].String).toID(),
			PnsID:               pgtype.Text{String: pnsID, Valid: true},
			Nip:                 pgtype.Text{String: nip, Valid: true},
		}); err != nil {
			return fmt.Errorf("repo create: %w", err)
		}

	case upd.ActionUpdate:
		if _, err := sqlcTx.UpdateRiwayatPendidikan(ctx, sqlc.UpdateRiwayatPendidikanParams{
			ID:                  id,
			Nip:                 nip,
			TingkatPendidikanID: data.TingkatPendidikanID[1],
			PendidikanID:        data.PendidikanID[1],
			NamaSekolah:         data.NamaSekolah[1],
			TahunLulus:          data.TahunLulus[1],
			NoIjazah:            data.NomorIjazah[1],
			GelarDepan:          data.GelarDepan[1],
			GelarBelakang:       data.GelarBelakang[1],
			NegaraSekolah:       data.NegaraSekolah[1],
			TugasBelajar:        statusBelajar(data.TugasBelajar[1].String).toID(),
		}); err != nil {
			return fmt.Errorf("repo update: %w", err)
		}

	case upd.ActionDelete:
		if _, err := sqlcTx.DeleteRiwayatPendidikan(ctx, sqlc.DeleteRiwayatPendidikanParams{
			ID:  id,
			Nip: nip,
		}); err != nil {
			return fmt.Errorf("repo delete: %w", err)
		}

	default:
		return errors.New("unimplemented action")
	}

	return nil
}
