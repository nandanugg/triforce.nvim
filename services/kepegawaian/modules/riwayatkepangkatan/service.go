package riwayatkepangkatan

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatKepangkatan(ctx context.Context, arg dbrepo.ListRiwayatKepangkatanParams) ([]dbrepo.ListRiwayatKepangkatanRow, error)
	CountRiwayatKepangkatan(ctx context.Context, pnsNip string) (int64, error)
	GetBerkasRiwayatKepangkatan(ctx context.Context, arg dbrepo.GetBerkasRiwayatKepangkatanParams) (pgtype.Text, error)
	GetPegawaiByNIP(ctx context.Context, nip string) (dbrepo.GetPegawaiByNIPRow, error)
	GetJenisKenaikanPangkat(ctx context.Context, id int32) (dbrepo.GetJenisKenaikanPangkatRow, error)
	GetRefGolongan(ctx context.Context, id int32) (dbrepo.GetRefGolonganRow, error)

	CreateRiwayatKepangkatan(ctx context.Context, arg dbrepo.CreateRiwayatKepangkatanParams) (string, error)
	UpdateRiwayatKepangkatan(ctx context.Context, arg dbrepo.UpdateRiwayatKepangkatanParams) (int64, error)
	DeleteRiwayatKepangkatan(ctx context.Context, arg dbrepo.DeleteRiwayatKepangkatanParams) (int64, error)
	UploadBerkasRiwayatKepangkatan(ctx context.Context, arg dbrepo.UploadBerkasRiwayatKepangkatanParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatKepangkatan, uint, error) {
	data, err := s.repo.ListRiwayatKepangkatan(ctx, dbrepo.ListRiwayatKepangkatanParams{
		PnsNip: nip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo ListRiwayatKepangkatan: %w", err)
	}

	count, err := s.repo.CountRiwayatKepangkatan(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo CountRiwayatKepangkatan: %w", err)
	}

	result := typeutil.Map(data, func(row dbrepo.ListRiwayatKepangkatanRow) riwayatKepangkatan {
		return riwayatKepangkatan{
			ID:                        row.ID,
			IDJenisKP:                 row.JenisKpID,
			NamaJenisKP:               row.NamaJenisKp.String,
			IDGolongan:                row.GolonganID,
			NamaGolongan:              row.NamaGolongan.String,
			NamaGolonganPangkat:       row.NamaGolonganPangkat.String,
			TMTGolongan:               db.Date(row.TmtGolongan.Time),
			SKNomor:                   row.SkNomor.String,
			SKTanggal:                 db.Date(row.SkTanggal.Time),
			MKGolonganTahun:           row.MkGolonganTahun,
			MKGolonganBulan:           row.MkGolonganBulan,
			NoBKN:                     row.NoBkn.String,
			TanggalBKN:                db.Date(row.TanggalBkn.Time),
			JumlahAngkaKreditTambahan: row.JumlahAngkaKreditTambahan,
			JumlahAngkaKreditUtama:    row.JumlahAngkaKreditUtama,
		}
	})

	return result, uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id string) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatKepangkatan(ctx, dbrepo.GetBerkasRiwayatKepangkatanParams{
		PnsNip: pgtype.Text{String: nip, Valid: true},
		ID:     id,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}
	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}

func (s *service) create(ctx context.Context, nip string, params upsertParams) (string, error) {
	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		return "", errPegawaiNotFound
	}

	ref, err := s.validateReferences(ctx, params)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateRiwayatKepangkatan(ctx, dbrepo.CreateRiwayatKepangkatanParams{
		JenisKpID:                 pgtype.Int4{Int32: typeutil.FromPtr(params.JenisKPID), Valid: params.JenisKPID != nil},
		JenisKp:                   ref.jenisKP.nama,
		GolonganID:                pgtype.Int2{Int16: params.GolonganID, Valid: true},
		GolonganNama:              ref.golongan.nama,
		PangkatNama:               ref.golongan.pangkat,
		TmtGolongan:               params.TMTGolongan.ToPgtypeDate(),
		SkNomor:                   pgtype.Text{String: params.NomorSK, Valid: true},
		SkTanggal:                 params.TanggalSK.ToPgtypeDate(),
		NoBkn:                     pgtype.Text{String: params.NomorBKN, Valid: params.NomorBKN != ""},
		TanggalBkn:                params.TanggalBKN.ToPgtypeDate(),
		MkGolonganTahun:           pgtype.Int2{Int16: params.MasaKerjaGolonganTahun, Valid: true},
		MkGolonganBulan:           pgtype.Int2{Int16: params.MasaKerjaGolonganBulan, Valid: true},
		JumlahAngkaKreditUtama:    pgtype.Int4{Int32: typeutil.FromPtr(params.JumlahAngkaKreditUtama), Valid: params.JumlahAngkaKreditUtama != nil},
		JumlahAngkaKreditTambahan: pgtype.Int4{Int32: typeutil.FromPtr(params.JumlahAngkaKreditTambahan), Valid: params.JumlahAngkaKreditTambahan != nil},
		PnsID:                     pgtype.Text{String: pegawai.PnsID, Valid: true},
		PnsNip:                    pgtype.Text{String: nip, Valid: true},
		PnsNama:                   pegawai.Nama,
	})
	if err != nil {
		return "", fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id, nip string, params upsertParams) (bool, error) {
	ref, err := s.validateReferences(ctx, params)
	if err != nil {
		return false, err
	}

	affected, err := s.repo.UpdateRiwayatKepangkatan(ctx, dbrepo.UpdateRiwayatKepangkatanParams{
		ID:                        id,
		Nip:                       nip,
		JenisKpID:                 pgtype.Int4{Int32: typeutil.FromPtr(params.JenisKPID), Valid: params.JenisKPID != nil},
		JenisKp:                   ref.jenisKP.nama,
		GolonganID:                pgtype.Int2{Int16: params.GolonganID, Valid: true},
		GolonganNama:              ref.golongan.nama,
		PangkatNama:               ref.golongan.pangkat,
		TmtGolongan:               params.TMTGolongan.ToPgtypeDate(),
		SkNomor:                   pgtype.Text{String: params.NomorSK, Valid: true},
		SkTanggal:                 params.TanggalSK.ToPgtypeDate(),
		NoBkn:                     pgtype.Text{String: params.NomorBKN, Valid: params.NomorBKN != ""},
		TanggalBkn:                params.TanggalBKN.ToPgtypeDate(),
		MkGolonganTahun:           pgtype.Int2{Int16: params.MasaKerjaGolonganTahun, Valid: true},
		MkGolonganBulan:           pgtype.Int2{Int16: params.MasaKerjaGolonganBulan, Valid: true},
		JumlahAngkaKreditUtama:    pgtype.Int4{Int32: typeutil.FromPtr(params.JumlahAngkaKreditUtama), Valid: params.JumlahAngkaKreditUtama != nil},
		JumlahAngkaKreditTambahan: pgtype.Int4{Int32: typeutil.FromPtr(params.JumlahAngkaKreditTambahan), Valid: params.JumlahAngkaKreditTambahan != nil},
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatKepangkatan(ctx, dbrepo.DeleteRiwayatKepangkatanParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatKepangkatan(ctx, dbrepo.UploadBerkasRiwayatKepangkatanParams{
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
	golongan struct {
		nama    pgtype.Text
		pangkat pgtype.Text
	}
	jenisKP struct {
		nama pgtype.Text
	}
}

func (s *service) validateReferences(ctx context.Context, params upsertParams) (*references, error) {
	var errs []error
	ref := references{}

	golongan, err := s.repo.GetRefGolongan(ctx, int32(params.GolonganID))
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get golongan: %w", err)
		}
		errs = append(errs, errGolonganNotFound)
	}
	ref.golongan.nama = golongan.Nama
	ref.golongan.pangkat = golongan.NamaPangkat

	if params.JenisKPID != nil {
		jenisKP, err := s.repo.GetJenisKenaikanPangkat(ctx, *params.JenisKPID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("repo get jenis kenaikan pangkat: %w", err)
			}
			errs = append(errs, errJenisKenaikanPangkatNotFound)
		}
		ref.jenisKP.nama = jenisKP.Nama
	}

	if len(errs) > 0 {
		return nil, api.NewMultiError(errs)
	}
	return &ref, nil
}
