package keluarga

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListOrangTuaByNip(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListOrangTuaByNipRow, error)
	ListPasanganByNip(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListPasanganByNipRow, error)
	ListAnakByNip(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListAnakByNipRow, error)
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
		return s.mapListOrangTua(row)
	})

	// 2. Pasangan
	pasanganList, err := s.repo.ListPasanganByNip(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list pasangan: %w", err)
	}
	result.Pasangan = typeutil.Map(pasanganList, func(row repo.ListPasanganByNipRow) pasangan {
		return s.mapListPasangan(row)
	})

	// 3. Anak
	anakList, err := s.repo.ListAnakByNip(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list anak: %w", err)
	}
	result.Anak = typeutil.Map(anakList, func(row repo.ListAnakByNipRow) anak {
		return s.mapListAnak(row)
	})

	return result, nil
}

func (s *service) mapListOrangTua(ot repo.ListOrangTuaByNipRow) orangTua {
	return orangTua{
		ID:          ot.ID,
		Hubungan:    hubunganToPeran(ot.Hubungan),
		StatusHidup: statusHidupFromTanggalMeninggal(ot.TanggalMeninggal),
		Nama:        ot.Nama.String,
		Agama:       ot.AgamaNama.String,
		NIK:         ot.Nik.String,
	}
}

func (s *service) mapListPasangan(p repo.ListPasanganByNipRow) pasangan {
	return pasangan{
		ID:               p.ID,
		StatusPNS:        pnsToLabel(p.Pns),
		Nama:             p.Nama.String,
		TanggalMenikah:   db.Date(p.TanggalMenikah.Time),
		Karsus:           p.Karsus.String,
		StatusNikah:      statusPernikahanToString(p.Status),
		Agama:            p.Agama.String,
		NIK:              "", // TODO: map actual nik if available
		AkteNikah:        p.AkteNikah.String,
		AkteMeninggal:    p.AkteMeninggal.String,
		AkteCerai:        p.AkteCerai.String,
		TanggalMeninggal: db.Date(p.TanggalMeninggal.Time),
		TanggalCerai:     db.Date(p.TanggalCerai.Time),
		TanggalLahir:     db.Date(p.TanggalLahir.Time),
	}
}

func (s *service) mapListAnak(a repo.ListAnakByNipRow) anak {
	return anak{
		ID:           a.ID,
		Nama:         a.Nama.String,
		NIK:          "", // TODO: map actual nik if available
		JenisKelamin: a.JenisKelamin.String,
		StatusAnak:   statusAnakToLabel(a.StatusAnak),
		NamaOrangTua: a.NamaOrangTua.String,
		TanggalLahir: db.Date(a.TanggalLahir.Time),
		AnakKe:       a.AnakKe,
	}
}
