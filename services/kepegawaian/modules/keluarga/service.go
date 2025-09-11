package keluarga

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

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
func (s *service) list(ctx context.Context, NIP string) (keluarga, error) {
	var result keluarga
	result.OrangTua = []orangTua{}
	result.Pasangan = []pasangan{}
	result.Anak = []anak{}

	// 1. Orang Tua
	orangTuaList, err := s.repo.ListOrangTuaByNip(ctx, pgtype.Text{String: NIP, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list orang tua: %w", err)
	}

	for _, orangTua := range orangTuaList {
		result.OrangTua = append(result.OrangTua, s.mapListOrangTua(orangTua))
	}

	// 2. Pasangan
	pasanganList, err := s.repo.ListPasanganByNip(ctx, pgtype.Text{String: NIP, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list pasangan: %w", err)
	}
	for _, pasangan := range pasanganList {
		result.Pasangan = append(result.Pasangan, s.mapListPasangan(pasangan))
	}

	// 3. Anak
	anakList, err := s.repo.ListAnakByNip(ctx, pgtype.Text{String: NIP, Valid: true})
	if err != nil {
		return keluarga{}, fmt.Errorf("repo list anak: %w", err)
	}
	for _, anak := range anakList {
		result.Anak = append(result.Anak, s.mapListAnak(anak))
	}

	return result, nil
}

func (s *service) mapListOrangTua(ot repo.ListOrangTuaByNipRow) orangTua {
	return orangTua{
		Hubungan:        HubunganToPeran(ot.Hubungan),
		StatusMeninggal: StatusHidupFromTanggalMeninggal(ot.TglMeninggal),
		Nama:            nullStringPtr(ot.Nama),
		Agama:           nullStringPtr(ot.AgamaNama),
		Nik:             nullStringPtr(ot.Nik),
	}
}

func (s *service) mapListPasangan(p repo.ListPasanganByNipRow) pasangan {
	return pasangan{
		Nama:           nullStringPtr(p.Nama),
		Nik:            nullStringPtr(p.Nik),
		StatusPNS:      PNSToLabel(p.Pns),
		Agama:          nullStringPtr(p.AgamaNama),
		StatusNikah:    StatusPernikahanToString(p.Status),
		TanggalMenikah: &p.TanggalMenikah.Time,
	}
}

func (s *service) mapListAnak(a repo.ListAnakByNipRow) anak {
	return anak{
		Nama:         nullStringPtr(a.Nama),
		Nip:          nullStringPtr(a.Nip),
		JenisKelamin: JenisKelaminToLabel(a.JenisKelamin),
		StatusAnak:   StatusAnakToLabel(a.StatusAnak),
		TanggalLahir: &a.TanggalLahir.Time,
		AnakKe:       &a.AnakKe,
	}
}
