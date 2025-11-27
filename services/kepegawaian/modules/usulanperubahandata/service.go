package usulanperubahandata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type ServiceInterface interface {
	GeneratePerubahanData(ctx context.Context, nip, action, dataID string, requestData json.RawMessage) ([]byte, error)
	SyncPerubahanData(ctx context.Context, sqlcTx *sqlc.Queries, nip, action, dataID string, perubahanData []byte) error
}

type repository interface {
	ListUnreadUsulanPerubahanDataByNIP(ctx context.Context, arg sqlc.ListUnreadUsulanPerubahanDataByNIPParams) ([]sqlc.ListUnreadUsulanPerubahanDataByNIPRow, error)
	CountUnreadUsulanPerubahanDataByNIP(ctx context.Context, arg sqlc.CountUnreadUsulanPerubahanDataByNIPParams) (int64, error)
	ListPendingUsulanPerubahanData(ctx context.Context, arg sqlc.ListPendingUsulanPerubahanDataParams) ([]sqlc.ListPendingUsulanPerubahanDataRow, error)
	CountPendingUsulanPerubahanData(ctx context.Context, arg sqlc.CountPendingUsulanPerubahanDataParams) (int64, error)
	GetUsulanPerubahanData(ctx context.Context, arg sqlc.GetUsulanPerubahanDataParams) (sqlc.GetUsulanPerubahanDataRow, error)
	ListUnitKerjaLengkapByIDs(ctx context.Context, ids []string) ([]sqlc.ListUnitKerjaLengkapByIDsRow, error)

	CreateUsulanPerubahanData(ctx context.Context, arg sqlc.CreateUsulanPerubahanDataParams) (int64, error)
	DeleteUsulanPerubahanData(ctx context.Context, arg sqlc.DeleteUsulanPerubahanDataParams) error
	MarkAsReadUsulanPerubahanData(ctx context.Context, arg sqlc.MarkAsReadUsulanPerubahanDataParams) error
	UpdateStatusUsulanPerubahanData(ctx context.Context, arg sqlc.UpdateStatusUsulanPerubahanDataParams) error

	WithTx(tx pgx.Tx) *sqlc.Queries
}

type service struct {
	db   *pgxpool.Pool
	repo repository
}

func newService(db *pgxpool.Pool, r repository) *service {
	return &service{db: db, repo: r}
}

func (s *service) listUnread(ctx context.Context, nip, jenisData string, limit, offset uint) ([]usulanPerubahanData, uint, error) {
	rows, err := s.repo.ListUnreadUsulanPerubahanDataByNIP(ctx, sqlc.ListUnreadUsulanPerubahanDataByNIPParams{
		Nip:       nip,
		JenisData: jenisData,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountUnreadUsulanPerubahanDataByNIP(ctx, sqlc.CountUnreadUsulanPerubahanDataByNIPParams{
		Nip:       nip,
		JenisData: jenisData,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	res := make([]usulanPerubahanData, 0, len(rows))
	for i, row := range rows {
		res = append(res, usulanPerubahanData{
			ID:        row.ID,
			JenisData: row.JenisData,
			DataID:    &rows[i].DataID,
			Action:    row.Action,
			Status:    row.Status,
			Catatan:   &rows[i].Catatan.String,
		})

		if err := json.Unmarshal(row.PerubahanData, &res[i].PerubahanData); err != nil {
			return nil, 0, fmt.Errorf("json unmarshal: %w", err)
		}
	}
	return res, uint(count), nil
}

type listPendingOptions struct {
	nama        string
	nip         string
	jenisData   string
	unitKerjaID string
	kodeJabatan string
	golonganID  *int16
}

func (s *service) listPending(ctx context.Context, limit, offset uint, opts listPendingOptions) ([]usulanPerubahanData, uint, error) {
	rows, err := s.repo.ListPendingUsulanPerubahanData(ctx, sqlc.ListPendingUsulanPerubahanDataParams{
		Nama:        opts.nama,
		Nip:         opts.nip,
		JenisData:   opts.jenisData,
		UnitKerjaID: opts.unitKerjaID,
		KodeJabatan: opts.kodeJabatan,
		GolonganID:  pgtype.Int2{Int16: typeutil.FromPtr(opts.golonganID), Valid: opts.golonganID != nil},
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountPendingUsulanPerubahanData(ctx, sqlc.CountPendingUsulanPerubahanDataParams{
		Nama:        opts.nama,
		Nip:         opts.nip,
		JenisData:   opts.jenisData,
		UnitKerjaID: opts.unitKerjaID,
		KodeJabatan: opts.kodeJabatan,
		GolonganID:  pgtype.Int2{Int16: typeutil.FromPtr(opts.golonganID), Valid: opts.golonganID != nil},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	uniqUnorIDs := typeutil.UniqMap(rows, func(row sqlc.ListPendingUsulanPerubahanDataRow, _ int) string {
		return row.UnorIDPegawai.String
	})

	listUnorLengkap, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, uniqUnorIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("repo list unor: %w", err)
	}

	unorLengkapByID := typeutil.SliceToMap(listUnorLengkap, func(unorLengkap sqlc.ListUnitKerjaLengkapByIDsRow) (string, string) {
		return unorLengkap.ID, unorLengkap.NamaUnorLengkap
	})

	data := typeutil.Map(rows, func(row sqlc.ListPendingUsulanPerubahanDataRow) usulanPerubahanData {
		return usulanPerubahanData{
			ID:        row.ID,
			JenisData: row.JenisData,
			CreatedAt: row.CreatedAt.Time.Format(time.RFC3339),
			Pegawai: &pegawai{
				NIP:           row.Nip,
				Nama:          row.NamaPegawai.String,
				GelarDepan:    row.GelarDepanPegawai.String,
				GelarBelakang: row.GelarBelakangPegawai.String,
				UnitKerja:     unorLengkapByID[row.UnorIDPegawai.String],
			},
		}
	})
	return data, uint(count), nil
}

func (s *service) detail(ctx context.Context, jenisData string, id int64) (*usulanPerubahanData, error) {
	row, err := s.repo.GetUsulanPerubahanData(ctx, sqlc.GetUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get: %w", err)
	}

	unor, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, []string{row.UnorIDPegawai.String})
	if err != nil {
		return nil, fmt.Errorf("repo list unor: %w", err)
	}

	res := &usulanPerubahanData{
		ID:        row.ID,
		JenisData: row.JenisData,
		Action:    row.Action,
		Status:    row.Status,
		DataID:    &row.DataID,
		CreatedAt: row.CreatedAt.Time.Format(time.RFC3339),
		Pegawai: &pegawai{
			NIP:           row.Nip,
			Nama:          row.NamaPegawai.String,
			GelarDepan:    row.GelarDepanPegawai.String,
			GelarBelakang: row.GelarBelakangPegawai.String,
			Golongan:      typeutil.ToPtr(typeutil.Cast[string](row.GolonganPegawai)),
			Jabatan:       &row.JabatanPegawai.String,
			Photo:         &row.FotoPegawai,
			StatusPNS:     typeutil.ToPtr(typeutil.Cast[string](row.StatusPnsPegawai)),
		},
	}
	if len(unor) > 0 {
		res.Pegawai.UnitKerja = unor[0].NamaUnorLengkap
	}
	if err := json.Unmarshal(row.PerubahanData, &res.PerubahanData); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return res, nil
}

func (s *service) create(ctx context.Context, svc ServiceInterface, nip, jenisData string, req createRequest) error {
	bytes, err := svc.GeneratePerubahanData(ctx, nip, req.Action, req.DataID, req.Data)
	if err != nil {
		return err
	}

	if _, err := s.repo.CreateUsulanPerubahanData(ctx, sqlc.CreateUsulanPerubahanDataParams{
		Nip:           nip,
		JenisData:     jenisData,
		DataID:        pgtype.Text{String: req.DataID, Valid: req.DataID != ""},
		PerubahanData: bytes,
		Action:        req.Action,
	}); err != nil {
		return fmt.Errorf("repo create: %w", err)
	}
	return nil
}

func (s *service) markAsRead(ctx context.Context, nip, jenisData string, id int64) (bool, error) {
	row, err := s.repo.GetUsulanPerubahanData(ctx, sqlc.GetUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	})
	if err != nil || row.Nip != nip {
		if row.Nip != nip || errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repo get: %w", err)
	}

	if row.Status == statusDiusulkan {
		return false, errInvalidStateTransition
	}

	if err := s.repo.MarkAsReadUsulanPerubahanData(ctx, sqlc.MarkAsReadUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	}); err != nil {
		return false, fmt.Errorf("repo mark as read: %w", err)
	}
	return true, nil
}

func (s *service) delete(ctx context.Context, nip, jenisData string, id int64) (bool, error) {
	row, err := s.repo.GetUsulanPerubahanData(ctx, sqlc.GetUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	})
	if err != nil || row.Nip != nip {
		if row.Nip != nip || errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repo get: %w", err)
	}

	if row.Status != statusDiusulkan {
		return false, errInvalidStateTransition
	}

	if err := s.repo.DeleteUsulanPerubahanData(ctx, sqlc.DeleteUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	}); err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}
	return true, nil
}

func (s *service) reject(ctx context.Context, jenisData string, id int64, catatan string) (bool, error) {
	row, err := s.repo.GetUsulanPerubahanData(ctx, sqlc.GetUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repo get: %w", err)
	}

	if row.Status != statusDiusulkan {
		return false, errInvalidStateTransition
	}

	if err := s.repo.UpdateStatusUsulanPerubahanData(ctx, sqlc.UpdateStatusUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
		Status:    statusDitolak,
		Catatan:   pgtype.Text{String: catatan, Valid: true},
	}); err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}
	return true, nil
}

func (s *service) approve(ctx context.Context, svc ServiceInterface, jenisData string, id int64) (bool, error) {
	row, err := s.repo.GetUsulanPerubahanData(ctx, sqlc.GetUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repo get: %w", err)
	}

	if row.Status != statusDiusulkan {
		return false, errInvalidStateTransition
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			slog.WarnContext(ctx, "Error rollback transaction in usulan perubahan data", "error", err)
		}
	}()

	sqlcTx := s.repo.WithTx(tx)
	if err := svc.SyncPerubahanData(ctx, sqlcTx, row.Nip, row.Action, row.DataID.String, row.PerubahanData); err != nil {
		return false, fmt.Errorf("sync perubahan data: %w", err)
	}

	if err := sqlcTx.UpdateStatusUsulanPerubahanData(ctx, sqlc.UpdateStatusUsulanPerubahanDataParams{
		ID:        id,
		JenisData: jenisData,
		Status:    statusDisetujui,
	}); err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit tx: %w", err)
	}
	return true, nil
}
