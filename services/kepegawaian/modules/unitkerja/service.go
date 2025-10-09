package unitkerja

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListAkarUnitKerja(ctx context.Context, arg repo.ListAkarUnitKerjaParams) ([]repo.ListAkarUnitKerjaRow, error)
	ListUnitKerjaByDiatasanID(ctx context.Context, arg repo.ListUnitKerjaByDiatasanIDParams) ([]repo.ListUnitKerjaByDiatasanIDRow, error)
	ListUnitKerjaByNamaOrInduk(ctx context.Context, arg repo.ListUnitKerjaByNamaOrIndukParams) ([]repo.ListUnitKerjaByNamaOrIndukRow, error)
	CountAkarUnitKerja(ctx context.Context) (int64, error)
	CountUnitKerjaByDiatasanID(ctx context.Context, diatasanID pgtype.Text) (int64, error)
	CountUnitKerja(ctx context.Context, arg repo.CountUnitKerjaParams) (int64, error)
	GetUnitKerja(ctx context.Context, id string) (repo.GetUnitKerjaRow, error)
	GetProfilePegawaiByPNSID(ctx context.Context, pnsID string) (repo.GetProfilePegawaiByPNSIDRow, error)
	CreateUnitKerja(ctx context.Context, arg repo.CreateUnitKerjaParams) (repo.CreateUnitKerjaRow, error)
	UpdateUnitKerja(ctx context.Context, arg repo.UpdateUnitKerjaParams) (repo.UpdateUnitKerjaRow, error)
	DeleteUnitKerja(ctx context.Context, id string) (int64, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

type listParams struct {
	nama      string
	unorInduk string
	limit     uint
	offset    uint
}

func (s *service) list(ctx context.Context, arg listParams) ([]unitKerjaPublic, int64, error) {
	rows, err := s.repo.ListUnitKerjaByNamaOrInduk(ctx, repo.ListUnitKerjaByNamaOrIndukParams{
		UnorInduk: arg.unorInduk,
		Limit:     int32(arg.limit),
		Offset:    int32(arg.offset),
		Nama:      arg.nama,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error getUnitKerjaByNamaOrInduk: %w", err)
	}

	result := typeutil.Map(rows, func(row repo.ListUnitKerjaByNamaOrIndukRow) unitKerjaPublic {
		return unitKerjaPublic{
			ID:   row.ID,
			Nama: row.NamaUnor.String,
		}
	})

	total, err := s.repo.CountUnitKerja(ctx, repo.CountUnitKerjaParams{
		Nama:      arg.nama,
		UnorInduk: arg.unorInduk,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error countUnitKerja: %w", err)
	}
	return result, total, nil
}

type listAkarParams struct {
	limit  uint
	offset uint
}

func (s *service) listAkar(ctx context.Context, arg listAkarParams) ([]anakUnitKerja, int64, error) {
	rows, err := s.repo.ListAkarUnitKerja(ctx, repo.ListAkarUnitKerjaParams{
		Limit:  int32(arg.limit),
		Offset: int32(arg.offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listAkar] error listAkarUnitKerja: %w", err)
	}

	result := typeutil.Map(rows, func(row repo.ListAkarUnitKerjaRow) anakUnitKerja {
		return anakUnitKerja{
			ID:      row.ID,
			Nama:    row.NamaUnor.String,
			HasAnak: row.HasAnak,
		}
	})

	total, err := s.repo.CountAkarUnitKerja(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[listAkar] error countAkarUnitKerja: %w", err)
	}
	return result, total, nil
}

type listAnakParams struct {
	limit  uint
	offset uint
	id     string
}

func (s *service) listAnak(ctx context.Context, arg listAnakParams) ([]anakUnitKerja, int64, error) {
	pgID := pgtype.Text{Valid: arg.id != "", String: arg.id}
	rows, err := s.repo.ListUnitKerjaByDiatasanID(ctx, repo.ListUnitKerjaByDiatasanIDParams{
		Limit:      int32(arg.limit),
		Offset:     int32(arg.offset),
		DiatasanID: pgID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listAnak] error listByDiatasanID: %w", err)
	}

	result := typeutil.Map(rows, func(row repo.ListUnitKerjaByDiatasanIDRow) anakUnitKerja {
		return anakUnitKerja{
			ID:      row.ID,
			Nama:    row.NamaUnor.String,
			HasAnak: row.HasAnak,
		}
	})

	total, err := s.repo.CountUnitKerjaByDiatasanID(ctx, pgID)
	if err != nil {
		return nil, 0, fmt.Errorf("[listAkar] error countUnitKerjaByDiatasanID: %w", err)
	}
	return result, total, nil
}

func (s *service) get(ctx context.Context, id string) (*unitKerjaWithInduk, error) {
	row, err := s.repo.GetUnitKerja(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error getUnitKerja: %w", err)
	}

	result := &unitKerjaWithInduk{
		ID:                   row.ID,
		No:                   row.No.Int32,
		KodeInternal:         row.KodeInternal.String,
		Nama:                 row.Nama.String,
		EselonID:             row.EselonID.String,
		CepatKode:            row.CepatKode.String,
		NamaJabatan:          row.NamaJabatan.String,
		NamaPejabat:          row.NamaPejabat.String,
		DiatasanID:           row.DiatasanID.String,
		InstansiID:           row.InstansiID.String,
		PemimpinPnsID:        row.PemimpinPnsID.String,
		JenisUnorID:          row.JenisUnorID.String,
		UnorInduk:            row.UnorInduk.String,
		JumlahIdealStaff:     row.JumlahIdealStaff.Int16,
		Order:                row.Order.Int32,
		IsSatker:             row.IsSatker,
		Eselon1:              row.Eselon1.String,
		Eselon2:              row.Eselon2.String,
		Eselon3:              row.Eselon3.String,
		Eselon4:              row.Eselon4.String,
		ExpiredDate:          db.Date(row.ExpiredDate.Time),
		Keterangan:           row.Keterangan.String,
		JenisSatker:          row.JenisSatker.String,
		Abbreviation:         row.Abbreviation.String,
		UnorIndukPenyetaraan: row.UnorIndukPenyetaraan.String,
		JabatanID:            row.JabatanID.String,
		Waktu:                row.Waktu.String,
		Peraturan:            row.Peraturan.String,
		Remark:               row.Remark.String,
		Aktif:                row.Aktif.Bool,
		EselonNama:           row.EselonNama.String,
		NamaDiatasan:         row.NamaDiatasan.String,
		NamaUnorInduk:        row.NamaUnorInduk.String,
	}

	return result, nil
}

type createParams struct {
	diatasanID    string
	id            string
	nama          string
	kodeInternal  string
	namaJabatan   string
	pemimpinPNSID string
	isSatker      bool
	unorInduk     string
	expiredDate   db.Date
	keterangan    string
	abbreviation  string
	waktu         string
	jenisSatker   string
	peraturan     string
}

func (s *service) create(ctx context.Context, params createParams) (*unitKerja, error) {
	var namaPejabat string

	if params.pemimpinPNSID != "" {
		pejabat, err := s.repo.GetProfilePegawaiByPNSID(ctx, params.pemimpinPNSID)
		if err != nil {
			return nil, fmt.Errorf("[create] error getProfilePegawaiByPNSID: %w", err)
		}

		namaPejabat = pejabat.Nama.String
	}

	row, err := s.repo.CreateUnitKerja(ctx, repo.CreateUnitKerjaParams{
		ID:            params.id,
		IsSatker:      params.isSatker,
		DiatasanID:    pgtype.Text{String: params.diatasanID, Valid: params.diatasanID != ""},
		Nama:          pgtype.Text{String: params.nama, Valid: params.nama != ""},
		KodeInternal:  pgtype.Text{String: params.kodeInternal, Valid: params.kodeInternal != ""},
		NamaJabatan:   pgtype.Text{String: params.namaJabatan, Valid: params.namaJabatan != ""},
		PemimpinPnsID: pgtype.Text{String: params.pemimpinPNSID, Valid: params.pemimpinPNSID != ""},
		NamaPejabat:   pgtype.Text{String: namaPejabat, Valid: namaPejabat != ""},
		UnorInduk:     pgtype.Text{String: params.unorInduk, Valid: params.unorInduk != ""},
		ExpiredDate:   pgtype.Date{Time: time.Time(params.expiredDate), Valid: !time.Time(params.expiredDate).IsZero()},
		Keterangan:    pgtype.Text{String: params.keterangan, Valid: params.keterangan != ""},
		Abbreviation:  pgtype.Text{String: params.abbreviation, Valid: params.abbreviation != ""},
		Waktu:         pgtype.Text{String: params.waktu, Valid: params.waktu != ""},
		JenisSatker:   pgtype.Text{String: params.jenisSatker, Valid: params.jenisSatker != ""},
		Peraturan:     pgtype.Text{String: params.peraturan, Valid: params.peraturan != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("[create] error createUnitKerja: %w", err)
	}

	result := &unitKerja{
		ID:                   row.ID,
		No:                   row.No.Int32,
		KodeInternal:         row.KodeInternal.String,
		Nama:                 row.Nama.String,
		EselonID:             row.EselonID.String,
		CepatKode:            row.CepatKode.String,
		NamaJabatan:          row.NamaJabatan.String,
		NamaPejabat:          row.NamaPejabat.String,
		DiatasanID:           row.DiatasanID.String,
		InstansiID:           row.InstansiID.String,
		PemimpinPnsID:        row.PemimpinPnsID.String,
		JenisUnorID:          row.JenisUnorID.String,
		UnorInduk:            row.UnorInduk.String,
		JumlahIdealStaff:     row.JumlahIdealStaff.Int16,
		Order:                row.Order.Int32,
		IsSatker:             row.IsSatker,
		Eselon1:              row.Eselon1.String,
		Eselon2:              row.Eselon2.String,
		Eselon3:              row.Eselon3.String,
		Eselon4:              row.Eselon4.String,
		ExpiredDate:          db.Date(row.ExpiredDate.Time),
		Keterangan:           row.Keterangan.String,
		JenisSatker:          row.JenisSatker.String,
		Abbreviation:         row.Abbreviation.String,
		UnorIndukPenyetaraan: row.UnorIndukPenyetaraan.String,
		JabatanID:            row.JabatanID.String,
		Waktu:                row.Waktu.String,
		Peraturan:            row.Peraturan.String,
		Remark:               row.Remark.String,
		Aktif:                row.Aktif.Bool,
		EselonNama:           row.EselonNama.String,
	}

	return result, nil
}

type updateParams struct {
	diatasanID    string
	id            string
	nama          string
	kodeInternal  string
	namaJabatan   string
	pemimpinPNSID string
	isSatker      bool
	unorInduk     string
	expiredDate   db.Date
	keterangan    string
	abbreviation  string
	waktu         string
	jenisSatker   string
	peraturan     string
}

func (s *service) update(ctx context.Context, params updateParams) (*unitKerja, error) {
	var namaPejabat string

	if params.pemimpinPNSID != "" {
		pejabat, err := s.repo.GetProfilePegawaiByPNSID(ctx, params.pemimpinPNSID)
		if err != nil {
			return nil, fmt.Errorf("[update] error getProfilePegawaiByPNSID: %w", err)
		}

		namaPejabat = pejabat.Nama.String
	}

	row, err := s.repo.UpdateUnitKerja(ctx, repo.UpdateUnitKerjaParams{
		ID:            params.id,
		IsSatker:      params.isSatker,
		DiatasanID:    pgtype.Text{String: params.diatasanID, Valid: params.diatasanID != ""},
		Nama:          pgtype.Text{String: params.nama, Valid: params.nama != ""},
		KodeInternal:  pgtype.Text{String: params.kodeInternal, Valid: params.kodeInternal != ""},
		NamaJabatan:   pgtype.Text{String: params.namaJabatan, Valid: params.namaJabatan != ""},
		PemimpinPnsID: pgtype.Text{String: params.pemimpinPNSID, Valid: params.pemimpinPNSID != ""},
		NamaPejabat:   pgtype.Text{String: namaPejabat, Valid: namaPejabat != ""},
		UnorInduk:     pgtype.Text{String: params.unorInduk, Valid: params.unorInduk != ""},
		ExpiredDate:   pgtype.Date{Time: time.Time(params.expiredDate), Valid: !time.Time(params.expiredDate).IsZero()},
		Keterangan:    pgtype.Text{String: params.keterangan, Valid: params.keterangan != ""},
		Abbreviation:  pgtype.Text{String: params.abbreviation, Valid: params.abbreviation != ""},
		Waktu:         pgtype.Text{String: params.waktu, Valid: params.waktu != ""},
		JenisSatker:   pgtype.Text{String: params.jenisSatker, Valid: params.jenisSatker != ""},
		Peraturan:     pgtype.Text{String: params.peraturan, Valid: params.peraturan != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error updateUnitKerja: %w", err)
	}

	result := &unitKerja{
		ID:                   row.ID,
		No:                   row.No.Int32,
		KodeInternal:         row.KodeInternal.String,
		Nama:                 row.Nama.String,
		EselonID:             row.EselonID.String,
		CepatKode:            row.CepatKode.String,
		NamaJabatan:          row.NamaJabatan.String,
		NamaPejabat:          row.NamaPejabat.String,
		DiatasanID:           row.DiatasanID.String,
		InstansiID:           row.InstansiID.String,
		PemimpinPnsID:        row.PemimpinPnsID.String,
		JenisUnorID:          row.JenisUnorID.String,
		UnorInduk:            row.UnorInduk.String,
		JumlahIdealStaff:     row.JumlahIdealStaff.Int16,
		Order:                row.Order.Int32,
		IsSatker:             row.IsSatker,
		Eselon1:              row.Eselon1.String,
		Eselon2:              row.Eselon2.String,
		Eselon3:              row.Eselon3.String,
		Eselon4:              row.Eselon4.String,
		ExpiredDate:          db.Date(row.ExpiredDate.Time),
		Keterangan:           row.Keterangan.String,
		JenisSatker:          row.JenisSatker.String,
		Abbreviation:         row.Abbreviation.String,
		UnorIndukPenyetaraan: row.UnorIndukPenyetaraan.String,
		JabatanID:            row.JabatanID.String,
		Waktu:                row.Waktu.String,
		Peraturan:            row.Peraturan.String,
		Remark:               row.Remark.String,
		Aktif:                row.Aktif.Bool,
		EselonNama:           row.EselonNama.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id string) (bool, error) {
	affected, err := s.repo.DeleteUnitKerja(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error deleteUnitKerja: %w", err)
	}
	return affected > 0, nil
}
