package datapribadi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetDataPribadi(ctx context.Context, arg sqlc.GetDataPribadiParams) (sqlc.GetDataPribadiRow, error)
	ListUnitKerjaHierarchy(ctx context.Context, id string) ([]sqlc.ListUnitKerjaHierarchyRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) get(ctx context.Context, nip string) (*dataPribadi, error) {
	data, err := s.repo.GetDataPribadi(ctx, sqlc.GetDataPribadiParams{
		NipBaru:                pgtype.Text{String: nip, Valid: true},
		JenisJabatanStruktural: jenisJabatanStruktural,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get data pribadi: %w", err)
	}

	unitOrganisasi := []string{}
	if data.UnorID.Valid {
		unitOrganisasiRows, err := s.repo.ListUnitKerjaHierarchy(ctx, data.UnorID.String)
		if err != nil {
			return nil, fmt.Errorf("get unit kerja hierarchy: %w", err)
		}

		unitOrganisasi = make([]string, 0, len(unitOrganisasiRows))
		for _, row := range unitOrganisasiRows {
			if row.NamaUnor.Valid {
				unitOrganisasi = append(unitOrganisasi, row.NamaUnor.String)
			}
		}
	}

	return &dataPribadi{
		Nama:                     data.Nama.String,
		GelarDepan:               data.GelarDepan.String,
		GelarBelakang:            data.GelarBelakang.String,
		JabatanAktual:            typeutil.Cast[string](data.JabatanAktual),
		JenisJabatanAktual:       typeutil.Cast[string](data.JenisJabatanAktual),
		TMTJabatan:               db.Date(data.TmtJabatan.Time),
		NIP:                      data.Nip.String,
		NIK:                      data.Nik.String,
		NomorKK:                  data.Kk.String,
		JenisKelamin:             data.JenisKelamin.String,
		TempatLahir:              data.TempatLahir.String,
		TanggalLahir:             db.Date(data.TglLahir.Time),
		TingkatPendidikan:        data.TingkatPendidikan.String,
		Pendidikan:               data.Pendidikan.String,
		StatusPerkawinan:         data.JenisKawin.String,
		Agama:                    data.Agama.String,
		EmailDikbud:              data.EmailDikbud.String,
		EmailPribadi:             data.Email.String,
		Alamat:                   data.Alamat.String,
		NomorHP:                  data.NoHp.String,
		NomorKontakDarurat:       data.NoDarurat.String,
		JenisPegawai:             data.JenisPegawai.String,
		MasaKerjaKeseluruhan:     typeutil.Cast[string](data.MasaKerjaKeseluruhan),
		MasaKerjaGolongan:        data.MasaKerjaGolongan.String,
		Jabatan:                  data.Jabatan.String,
		KelasJabatan:             data.KelasJabatan.String,
		LokasiKerja:              data.LokasiKerja.String,
		GolonganRuangAwal:        typeutil.Cast[string](data.GolonganAwal),
		GolonganRuangAkhir:       typeutil.Cast[string](data.GolonganAkhir),
		PangkatAkhir:             data.PangkatAkhir.String,
		TMTGolongan:              db.Date(data.TmtGolongan.Time),
		TMTASN:                   db.Date(data.TmtAsn.Time),
		NomorSKASN:               data.NoSkAsn.String,
		StatusASN:                typeutil.Cast[string](data.StatusAsn),
		TMTPNS:                   db.Date(typeutil.Cast[time.Time](data.TmtPns)),
		KartuPegawai:             data.KartuPegawai.String,
		NomorSuratDokter:         data.NoSuratDokter.String,
		TanggalSuratDokter:       db.Date(data.TglSuratDokter.Time),
		NomorSuratBebasNarkoba:   data.NoBebasNarkoba.String,
		TanggalSuratBebasNarkoba: db.Date(data.TglBebasNarkoba.Time),
		NomorCatatanPolisi:       data.NoCatatanPolisi.String,
		TanggalCatatanPolisi:     db.Date(data.TglCatatanPolisi.Time),
		AkteKelahiran:            data.AkteKelahiran.String,
		NomorBPJS:                data.Bpjs.String,
		NPWP:                     data.Npwp.String,
		TanggalNPWP:              db.Date(data.TglNpwp.Time),
		NomorTaspen:              data.NoTaspen.String,
		UnitOrganisasi:           unitOrganisasi,
	}, nil
}
