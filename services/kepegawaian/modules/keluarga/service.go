package keluarga

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetRefAgama(ctx context.Context, id int32) (repo.GetRefAgamaRow, error)
	GetRefJenisKawin(ctx context.Context, id int32) (repo.GetRefJenisKawinRow, error)
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	IsPasanganExistsByIDAndNIP(ctx context.Context, arg repo.IsPasanganExistsByIDAndNIPParams) (bool, error)

	ListOrangTuaByNip(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListOrangTuaByNipRow, error)
	CreateOrangTua(ctx context.Context, arg repo.CreateOrangTuaParams) (int32, error)
	UpdateOrangTua(ctx context.Context, arg repo.UpdateOrangTuaParams) (int64, error)
	DeleteOrangTua(ctx context.Context, arg repo.DeleteOrangTuaParams) (int64, error)

	ListPasanganByNip(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListPasanganByNipRow, error)
	CreatePasangan(ctx context.Context, arg repo.CreatePasanganParams) (int64, error)
	UpdatePasangan(ctx context.Context, arg repo.UpdatePasanganParams) (int64, error)
	DeletePasangan(ctx context.Context, arg repo.DeletePasanganParams) (int64, error)

	ListAnakByNip(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListAnakByNipRow, error)
	CreateAnak(ctx context.Context, arg repo.CreateAnakParams) (int64, error)
	UpdateAnak(ctx context.Context, arg repo.UpdateAnakParams) (int64, error)
	DeleteAnak(ctx context.Context, arg repo.DeleteAnakParams) (int64, error)
}
type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

// list returns a structured keluarga containing all family members categorized by type.
func (s *service) list(ctx context.Context, nip string) (keluarga, error) {
	var result keluarga

	// 1. Orang Tua
	orangTuaList, err := s.repo.ListOrangTuaByNip(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list orang tua: %w", err)
	}
	result.OrangTua = typeutil.Map(orangTuaList, func(row repo.ListOrangTuaByNipRow) orangTua {
		return orangTua{
			ID:               row.ID,
			Hubungan:         labelHubunganOrangTua[row.Hubungan.Int16],
			TanggalMeninggal: db.Date(row.TanggalMeninggal.Time),
			AkteMeninggal:    row.AkteMeninggal.String,
			StatusHidup:      statusHidupFromTanggalMeninggal(row.TanggalMeninggal),
			Nama:             row.Nama.String,
			AgamaID:          row.AgamaID,
			Agama:            row.Agama.String,
			NIK:              row.Nik.String,
		}
	})

	// 2. Pasangan
	pasanganList, err := s.repo.ListPasanganByNip(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list pasangan: %w", err)
	}
	result.Pasangan = typeutil.Map(pasanganList, func(row repo.ListPasanganByNipRow) pasangan {
		return pasangan{
			ID:                 row.ID,
			StatusPNS:          pnsToLabel(row.Pns),
			Nama:               row.Nama.String,
			TanggalMenikah:     db.Date(row.TanggalMenikah.Time),
			Karsus:             row.Karsus.String,
			StatusPernikahanID: row.StatusPernikahanID,
			StatusNikah:        row.StatusPernikahan.String,
			AgamaID:            row.AgamaID,
			Agama:              row.Agama.String,
			NIK:                row.Nik.String,
			AkteNikah:          row.AkteNikah.String,
			AkteMeninggal:      row.AkteMeninggal.String,
			AkteCerai:          row.AkteCerai.String,
			TanggalMeninggal:   db.Date(row.TanggalMeninggal.Time),
			TanggalCerai:       db.Date(row.TanggalCerai.Time),
			TanggalLahir:       db.Date(row.TanggalLahir.Time),
		}
	})

	// 3. Anak
	anakList, err := s.repo.ListAnakByNip(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list anak: %w", err)
	}
	result.Anak = typeutil.Map(anakList, func(row repo.ListAnakByNipRow) anak {
		return anak{
			ID:                 row.ID,
			Nama:               row.Nama.String,
			NIK:                row.Nik.String,
			JenisKelamin:       row.JenisKelamin.String,
			StatusAnak:         labelStatusAnak[row.StatusAnak.String],
			StatusSekolah:      labelStatusSekolah[row.StatusSekolah.Int16],
			AgamaID:            row.AgamaID,
			Agama:              row.Agama.String,
			StatusPernikahanID: row.JenisKawinID,
			StatusPernikahan:   row.StatusPernikahan.String,
			NamaOrangTua:       row.NamaOrangTua.String,
			PasanganOrangTuaID: row.PasanganID,
			TanggalLahir:       db.Date(row.TanggalLahir.Time),
			AnakKe:             row.AnakKe,
		}
	})

	return result, nil
}

func (s *service) createOrangTua(ctx context.Context, nip string, params upsertOrangTuaParams) (int32, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	if err := s.validateOrangTuaReferences(ctx, params); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateOrangTua(ctx, repo.CreateOrangTuaParams{
		Nama:             pgtype.Text{String: params.Nama, Valid: true},
		JenisDokumen:     pgtype.Text{String: "KTP", Valid: params.NIK != ""},
		NoDokumen:        pgtype.Text{String: params.NIK, Valid: params.NIK != ""},
		Hubungan:         params.Hubungan.toID(),
		AgamaID:          pgtype.Int2{Int16: typeutil.FromPtr(params.AgamaID), Valid: params.AgamaID != nil},
		TanggalMeninggal: params.TanggalMeninggal.ToPgtypeDate(),
		AkteMeninggal:    pgtype.Text{String: params.AkteMeninggal, Valid: params.AkteMeninggal != ""},
		PnsID:            pgtype.Text{String: pnsID, Valid: true},
		Nip:              pgtype.Text{String: nip, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create orang tua: %w", err)
	}

	return id, nil
}

func (s *service) updateOrangTua(ctx context.Context, id int32, nip string, params upsertOrangTuaParams) (bool, error) {
	if err := s.validateOrangTuaReferences(ctx, params); err != nil {
		return false, err
	}

	affected, err := s.repo.UpdateOrangTua(ctx, repo.UpdateOrangTuaParams{
		ID:               id,
		Nip:              nip,
		Nama:             pgtype.Text{String: params.Nama, Valid: true},
		JenisDokumen:     pgtype.Text{String: "KTP", Valid: params.NIK != ""},
		NoDokumen:        pgtype.Text{String: params.NIK, Valid: params.NIK != ""},
		Hubungan:         params.Hubungan.toID(),
		AgamaID:          pgtype.Int2{Int16: typeutil.FromPtr(params.AgamaID), Valid: params.AgamaID != nil},
		TanggalMeninggal: params.TanggalMeninggal.ToPgtypeDate(),
		AkteMeninggal:    pgtype.Text{String: params.AkteMeninggal, Valid: params.AkteMeninggal != ""},
	})
	if err != nil {
		return false, fmt.Errorf("repo update orang tua: %w", err)
	}

	return affected > 0, nil
}

func (s *service) deleteOrangTua(ctx context.Context, id int32, nip string) (bool, error) {
	affected, err := s.repo.DeleteOrangTua(ctx, repo.DeleteOrangTuaParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete orang tua: %w", err)
	}

	return affected > 0, nil
}

func (s *service) validateOrangTuaReferences(ctx context.Context, params upsertOrangTuaParams) error {
	var errs []error
	if params.AgamaID != nil {
		if _, err := s.repo.GetRefAgama(ctx, int32(*params.AgamaID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get agama: %w", err)
			}
			errs = append(errs, errAgamaNotFound)
		}
	}

	if len(errs) > 0 {
		return api.NewMultiError(errs)
	}
	return nil
}

func (s *service) createPasangan(ctx context.Context, nip string, params upsertPasanganParams) (int64, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	if err := s.validatePasanganReferences(ctx, params); err != nil {
		return 0, err
	}

	id, err := s.repo.CreatePasangan(ctx, repo.CreatePasanganParams{
		Nama:             pgtype.Text{String: params.Nama, Valid: true},
		Nik:              pgtype.Text{String: params.NIK, Valid: params.NIK != ""},
		Pns:              statusPNS(params.IsPNS),
		TanggalLahir:     params.TanggalLahir.ToPgtypeDate(),
		Karsus:           pgtype.Text{String: params.NoKarsus, Valid: params.NoKarsus != ""},
		AgamaID:          pgtype.Int2{Int16: typeutil.FromPtr(params.AgamaID), Valid: params.AgamaID != nil},
		Status:           pgtype.Int2{Int16: params.StatusPernikahanID, Valid: true},
		Hubungan:         params.Hubungan.toID(),
		TanggalMenikah:   params.TanggalMenikah.ToPgtypeDate(),
		AkteNikah:        pgtype.Text{String: params.AkteNikah, Valid: params.AkteNikah != ""},
		TanggalMeninggal: params.TanggalMeninggal.ToPgtypeDate(),
		AkteMeninggal:    pgtype.Text{String: params.AkteMeninggal, Valid: params.AkteMeninggal != ""},
		TanggalCerai:     params.TanggalCerai.ToPgtypeDate(),
		AkteCerai:        pgtype.Text{String: params.AkteCerai, Valid: params.AkteCerai != ""},
		PnsID:            pgtype.Text{String: pnsID, Valid: true},
		Nip:              pgtype.Text{String: nip, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create pasangan: %w", err)
	}

	return id, nil
}

func (s *service) updatePasangan(ctx context.Context, id int64, nip string, params upsertPasanganParams) (bool, error) {
	if err := s.validatePasanganReferences(ctx, params); err != nil {
		return false, err
	}

	affected, err := s.repo.UpdatePasangan(ctx, repo.UpdatePasanganParams{
		ID:               id,
		Nip:              nip,
		Nama:             pgtype.Text{String: params.Nama, Valid: true},
		Nik:              pgtype.Text{String: params.NIK, Valid: params.NIK != ""},
		Pns:              statusPNS(params.IsPNS),
		TanggalLahir:     params.TanggalLahir.ToPgtypeDate(),
		Karsus:           pgtype.Text{String: params.NoKarsus, Valid: params.NoKarsus != ""},
		AgamaID:          pgtype.Int2{Int16: typeutil.FromPtr(params.AgamaID), Valid: params.AgamaID != nil},
		Status:           pgtype.Int2{Int16: params.StatusPernikahanID, Valid: true},
		Hubungan:         params.Hubungan.toID(),
		TanggalMenikah:   params.TanggalMenikah.ToPgtypeDate(),
		AkteNikah:        pgtype.Text{String: params.AkteNikah, Valid: params.AkteNikah != ""},
		TanggalMeninggal: params.TanggalMeninggal.ToPgtypeDate(),
		AkteMeninggal:    pgtype.Text{String: params.AkteMeninggal, Valid: params.AkteMeninggal != ""},
		TanggalCerai:     params.TanggalCerai.ToPgtypeDate(),
		AkteCerai:        pgtype.Text{String: params.AkteCerai, Valid: params.AkteCerai != ""},
	})
	if err != nil {
		return false, fmt.Errorf("repo update pasangan: %w", err)
	}

	return affected > 0, nil
}

func (s *service) deletePasangan(ctx context.Context, id int64, nip string) (bool, error) {
	affected, err := s.repo.DeletePasangan(ctx, repo.DeletePasanganParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete pasangan: %w", err)
	}

	return affected > 0, nil
}

func (s *service) validatePasanganReferences(ctx context.Context, params upsertPasanganParams) error {
	var errs []error
	if params.AgamaID != nil {
		if _, err := s.repo.GetRefAgama(ctx, int32(*params.AgamaID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get agama: %w", err)
			}
			errs = append(errs, errAgamaNotFound)
		}
	}

	if _, err := s.repo.GetRefJenisKawin(ctx, int32(params.StatusPernikahanID)); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("repo get status pernikahan: %w", err)
		}
		errs = append(errs, errStatusPernikahanNotFound)
	}

	if len(errs) > 0 {
		return api.NewMultiError(errs)
	}
	return nil
}

func (s *service) createAnak(ctx context.Context, nip string, params upsertAnakParams) (int64, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	if err := s.validateAnakReferences(ctx, nip, params); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateAnak(ctx, repo.CreateAnakParams{
		Nama:          pgtype.Text{String: params.Nama, Valid: true},
		Nik:           pgtype.Text{String: params.NIK, Valid: params.NIK != ""},
		PasanganID:    pgtype.Int8{Int64: params.PasanganOrangTuaID, Valid: true},
		JenisKelamin:  pgtype.Text{String: params.JenisKelamin, Valid: true},
		TanggalLahir:  params.TanggalLahir.ToPgtypeDate(),
		AgamaID:       pgtype.Int2{Int16: typeutil.FromPtr(params.AgamaID), Valid: params.AgamaID != nil},
		JenisKawinID:  pgtype.Int2{Int16: params.StatusPernikahanID, Valid: true},
		StatusAnak:    params.StatusAnak.toID(),
		StatusSekolah: params.StatusSekolah.toID(),
		AnakKe:        pgtype.Int2{Int16: typeutil.FromPtr(params.AnakKe), Valid: typeutil.FromPtr(params.AnakKe) != 0},
		PnsID:         pgtype.Text{String: pnsID, Valid: true},
		Nip:           pgtype.Text{String: nip, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create anak: %w", err)
	}

	return id, nil
}

func (s *service) updateAnak(ctx context.Context, id int64, nip string, params upsertAnakParams) (bool, error) {
	if err := s.validateAnakReferences(ctx, nip, params); err != nil {
		return false, err
	}

	affected, err := s.repo.UpdateAnak(ctx, repo.UpdateAnakParams{
		ID:            id,
		Nip:           nip,
		Nama:          pgtype.Text{String: params.Nama, Valid: true},
		Nik:           pgtype.Text{String: params.NIK, Valid: params.NIK != ""},
		PasanganID:    pgtype.Int8{Int64: params.PasanganOrangTuaID, Valid: true},
		JenisKelamin:  pgtype.Text{String: params.JenisKelamin, Valid: true},
		TanggalLahir:  params.TanggalLahir.ToPgtypeDate(),
		AgamaID:       pgtype.Int2{Int16: typeutil.FromPtr(params.AgamaID), Valid: params.AgamaID != nil},
		JenisKawinID:  pgtype.Int2{Int16: params.StatusPernikahanID, Valid: true},
		StatusAnak:    params.StatusAnak.toID(),
		StatusSekolah: params.StatusSekolah.toID(),
		AnakKe:        pgtype.Int2{Int16: typeutil.FromPtr(params.AnakKe), Valid: typeutil.FromPtr(params.AnakKe) != 0},
	})
	if err != nil {
		return false, fmt.Errorf("repo update anak: %w", err)
	}

	return affected > 0, nil
}

func (s *service) deleteAnak(ctx context.Context, id int64, nip string) (bool, error) {
	affected, err := s.repo.DeleteAnak(ctx, repo.DeleteAnakParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete anak: %w", err)
	}

	return affected > 0, nil
}

func (s *service) validateAnakReferences(ctx context.Context, nip string, params upsertAnakParams) error {
	var errs []error
	if params.AgamaID != nil {
		if _, err := s.repo.GetRefAgama(ctx, int32(*params.AgamaID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get agama: %w", err)
			}
			errs = append(errs, errAgamaNotFound)
		}
	}

	if _, err := s.repo.GetRefJenisKawin(ctx, int32(params.StatusPernikahanID)); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("repo get status pernikahan: %w", err)
		}
		errs = append(errs, errStatusPernikahanNotFound)
	}

	if ok, err := s.repo.IsPasanganExistsByIDAndNIP(ctx, repo.IsPasanganExistsByIDAndNIPParams{
		ID:  params.PasanganOrangTuaID,
		Nip: nip,
	}); err != nil {
		return fmt.Errorf("repo pasangan exists: %w", err)
	} else if !ok {
		errs = append(errs, errPasanganOrangTuaNotFound)
	}

	if len(errs) > 0 {
		return api.NewMultiError(errs)
	}
	return nil
}
