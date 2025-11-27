package pegawai

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetProfilePegawaiByPNSID(ctx context.Context, pnsID string) (sqlc.GetProfilePegawaiByPNSIDRow, error)
	ListUnitKerjaHierarchy(ctx context.Context, id string) ([]sqlc.ListUnitKerjaHierarchyRow, error)
	ListPegawaiAktif(ctx context.Context, arg sqlc.ListPegawaiAktifParams) ([]sqlc.ListPegawaiAktifRow, error)
	ListUnitKerjaLengkapByIDs(ctx context.Context, ids []string) ([]sqlc.ListUnitKerjaLengkapByIDsRow, error)
	CountPegawaiAktif(ctx context.Context, arg sqlc.CountPegawaiAktifParams) (int64, error)
	GetDataPribadi(ctx context.Context, arg sqlc.GetDataPribadiParams) (sqlc.GetDataPribadiRow, error)
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	GetRefAgama(ctx context.Context, id int32) (sqlc.GetRefAgamaRow, error)
	GetRefGolongan(ctx context.Context, id int32) (sqlc.GetRefGolonganRow, error)
	GetRefJenisKawin(ctx context.Context, id int32) (sqlc.GetRefJenisKawinRow, error)
	GetRefTingkatPendidikan(ctx context.Context, id int32) (sqlc.GetRefTingkatPendidikanRow, error)
	GetRefJenisJabatan(ctx context.Context, id int32) (sqlc.GetRefJenisJabatanRow, error)
	WithTx(tx pgx.Tx) *sqlc.Queries
	UpdateDataPegawai(ctx context.Context, arg sqlc.UpdateDataPegawaiParams) error
	UpdateTTDPegawaiNIPByNIP(ctx context.Context, arg sqlc.UpdateTTDPegawaiNIPByNIPParams) error
	UpdateAnakNIPByNIP(ctx context.Context, arg sqlc.UpdateAnakNIPByNIPParams) error
	UpdateOrangTuaNIPByNIP(ctx context.Context, arg sqlc.UpdateOrangTuaNIPByNIPParams) error
	UpdatePasanganNIPByNIP(ctx context.Context, arg sqlc.UpdatePasanganNIPByNIPParams) error
	UpdateRiwayatAsesmenNineBoxNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatAsesmenNineBoxNamaNipByPNSIDParams) error
	UpdateRiwayatAsesmenNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatAsesmenNamaNipByPNSIDParams) error
	UpdateRiwayatHukumanDisiplinNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatHukumanDisiplinNamaNipByPNSIDParams) error
	UpdateRiwayatJabatanNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatJabatanNamaNipByPNSIDParams) error
	UpdateRiwayatKenaikanGajiBerkalaNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatKenaikanGajiBerkalaNamaNipByPNSIDParams) error
	UpdateRiwayatKepangkatanNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatKepangkatanNamaNipByPNSIDParams) error
	UpdateRiwayatKinerjaNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatKinerjaNamaNipByPNSIDParams) error
	UpdateRiwayatPelatihanFungsionalNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatPelatihanFungsionalNamaNipByNIPParams) error
	UpdateRiwayatPelatihanSIASNNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatPelatihanSIASNNamaNipByPNSIDParams) error
	UpdateRiwayatPelatihanStrukturalNamaNipByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatPelatihanStrukturalNamaNipByPNSIDParams) error
	UpdateRiwayatPelatihanTeknisNIPByPNSID(ctx context.Context, arg sqlc.UpdateRiwayatPelatihanTeknisNIPByPNSIDParams) error
	UpdateRiwayatPendidikanNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatPendidikanNamaNipByNIPParams) error
	UpdateRiwayatPenghargaanNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatPenghargaanNamaNipByNIPParams) error
	UpdateRiwayatPenugasanNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatPenugasanNamaNipByNIPParams) error
	UpdateRiwayatPindahUnitKerjaNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatPindahUnitKerjaNamaNipByNIPParams) error
	UpdateRiwayatSertifikasiNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatSertifikasiNamaNipByNIPParams) error
	UpdateRiwayatUjiKompetensiNamaNipByNIP(ctx context.Context, arg sqlc.UpdateRiwayatUjiKompetensiNamaNipByNIPParams) error
	UpdateRiwayatSuratKeputusanNamaNipPemrosesByNIP(ctx context.Context, arg sqlc.UpdateRiwayatSuratKeputusanNamaNipPemrosesByNIPParams) error
	UpdateSuratKeputusanNamaNipPemilikByNIP(ctx context.Context, arg sqlc.UpdateSuratKeputusanNamaNipPemilikByNIPParams) error
	UpdateSuratKeputusanNipPemrosesByNIP(ctx context.Context, arg sqlc.UpdateSuratKeputusanNipPemrosesByNIPParams) error
}

type service struct {
	repo repository
	db   *pgxpool.Pool
}

func newService(r repository, db *pgxpool.Pool) *service {
	return &service{repo: r, db: db}
}

func (s *service) getProfileByPNSID(ctx context.Context, pnsID string) (*profile, error) {
	data, err := s.repo.GetProfilePegawaiByPNSID(ctx, pnsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get profile: %w", err)
	}

	unitOrganisasi := make([]string, 0)
	if data.UnorID.Valid {
		rows, err := s.repo.ListUnitKerjaHierarchy(ctx, data.UnorID.String)
		if err != nil {
			return nil, fmt.Errorf("repo list unit kerja hierarchy: %w", err)
		}

		unitOrganisasi = typeutil.FilterMap(rows, func(row sqlc.ListUnitKerjaHierarchyRow) (string, bool) {
			return row.NamaUnor.String, row.NamaUnor.String != ""
		})
	}

	return &profile{
		NIPLama:        data.NipLama.String,
		NIPBaru:        data.NipBaru.String,
		GelarDepan:     data.GelarDepan.String,
		GelarBelakang:  data.GelarBelakang.String,
		Nama:           data.Nama.String,
		Pangkat:        data.Pangkat.String,
		Golongan:       typeutil.Cast[string](data.Golongan),
		Jabatan:        data.Jabatan.String,
		UnitOrganisasi: unitOrganisasi,
		Photo:          data.Foto,
	}, nil
}

type adminListPegawaiParams struct {
	limit      uint
	offset     uint
	keyword    string
	unitID     string
	golonganID int32
	jabatanID  string
	status     string
}

func (s *service) adminListPegawai(ctx context.Context, arg adminListPegawaiParams) ([]pegawai, uint, error) {
	statusHukum := getStatusHukum(arg.status)
	data, err := s.repo.ListPegawaiAktif(ctx, sqlc.ListPegawaiAktifParams{
		Limit:       int32(arg.limit),
		Offset:      int32(arg.offset),
		Keyword:     pgtype.Text{Valid: arg.keyword != "", String: arg.keyword},
		UnitKerjaID: pgtype.Text{Valid: arg.unitID != "", String: arg.unitID},
		GolonganID:  pgtype.Int4{Valid: arg.golonganID != 0, Int32: arg.golonganID},
		JabatanID:   pgtype.Text{Valid: arg.jabatanID != "", String: arg.jabatanID},
		StatusHukum: pgtype.Text{Valid: statusHukum != "", String: statusHukum},
		StatusPns:   getStatusPNSDB(arg.status),
		Mpp:         statusKedudukanHukumMPP,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[pegawai-adminListPegawai] repo ListPegawaiAktif: %w", err)
	}

	uniqUnorIDs := typeutil.UniqMap(data, func(row sqlc.ListPegawaiAktifRow, _ int) string {
		return row.UnorID.String
	})

	listUnorLengkap, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, uniqUnorIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("[pegawai-adminListPegawai] repo ListUnitKerjaLengkapByIDs: %w", err)
	}
	unorLengkapByID := typeutil.SliceToMap(listUnorLengkap, func(unorLengkap sqlc.ListUnitKerjaLengkapByIDsRow) (string, string) {
		return unorLengkap.ID, unorLengkap.NamaUnorLengkap
	})

	count, err := s.repo.CountPegawaiAktif(ctx, sqlc.CountPegawaiAktifParams{
		Keyword:     pgtype.Text{Valid: arg.keyword != "", String: arg.keyword},
		UnitKerjaID: pgtype.Text{Valid: arg.unitID != "", String: arg.unitID},
		GolonganID:  pgtype.Int4{Valid: arg.golonganID != 0, Int32: arg.golonganID},
		JabatanID:   pgtype.Text{Valid: arg.jabatanID != "", String: arg.jabatanID},
		StatusHukum: pgtype.Text{Valid: statusHukum != "", String: statusHukum},
		StatusPns:   getStatusPNSDB(arg.status),
		Mpp:         statusKedudukanHukumMPP,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[pegawai-adminListPegawai] repo CountPegawaiAktif: %w", err)
	}

	result := typeutil.Map(data, func(row sqlc.ListPegawaiAktifRow) pegawai {
		return pegawai{
			PNSID:         row.PnsID,
			NIP:           row.Nip.String,
			GelarDepan:    row.GelarDepan.String,
			Nama:          row.Nama.String,
			GelarBelakang: row.GelarBelakang.String,
			Golongan:      row.Golongan.String,
			Jabatan:       row.Jabatan.String,
			UnitKerja:     unorLengkapByID[row.UnorID.String],
			Status:        getLabelStatusPNS(row.NamaKedudukuanHukum.String, row.StatusCpnsPns.String),
			Photo:         row.Foto,
		}
	})

	return result, uint(count), nil
}

func (s *service) adminGetPegawai(ctx context.Context, nip string) (*pegawaiDetail, error) {
	data, err := s.repo.GetDataPribadi(ctx, sqlc.GetDataPribadiParams{
		NipBaru:                pgtype.Text{String: nip, Valid: true},
		JenisJabatanStruktural: jenisJabatanStruktural,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[pegawai-adminGetPegawai] repo get data pribadi: %w", err)
	}

	var (
		tmtPNS         db.Date
		statusPNS      string
		unitOrganisasi = make([]string, 0)
	)
	if data.UnorID.Valid {
		rows, err := s.repo.ListUnitKerjaHierarchy(ctx, data.UnorID.String)
		if err != nil {
			return nil, fmt.Errorf("[pegawai-adminGetPegawai] repo list unit kerja hierarchy: %w", err)
		}

		unitOrganisasi = typeutil.FilterMap(rows, func(row sqlc.ListUnitKerjaHierarchyRow) (string, bool) {
			return row.NamaUnor.String, row.NamaUnor.String != ""
		})
	}
	if data.StatusPns.String == "P" || data.StatusPns.String == "PNS" {
		statusPNS, tmtPNS = "PNS", db.Date(data.TmtPns.Time)
	}

	return &pegawaiDetail{
		PNSID:                    data.PnsID,
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
		TanggalLahir:             db.Date(data.TanggalLahir.Time),
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
		JenisJabatan:             data.JenisJabatan.String,
		KelasJabatan:             data.KelasJabatan.String,
		LokasiKerja:              data.LokasiKerja.String,
		GolonganRuangAwal:        typeutil.Cast[string](data.GolonganAwal),
		GolonganRuangAkhir:       typeutil.Cast[string](data.GolonganAkhir),
		PangkatAkhir:             data.PangkatAkhir.String,
		TMTGolongan:              db.Date(data.TmtGolongan.Time),
		TMTASN:                   db.Date(data.TmtAsn.Time),
		NomorSKASN:               data.NoSkAsn.String,
		IsPPPK:                   data.IsPppk.Bool,
		StatusASN:                data.StatusAsn.String,
		StatusPNS:                statusPNS,
		TMTPNS:                   tmtPNS,
		KartuPegawai:             data.KartuPegawai.String,
		NomorSuratDokter:         data.NoSuratDokter.String,
		TanggalSuratDokter:       db.Date(data.TanggalSuratDokter.Time),
		NomorSuratBebasNarkoba:   data.NoBebasNarkoba.String,
		TanggalSuratBebasNarkoba: db.Date(data.TanggalBebasNarkoba.Time),
		NomorCatatanPolisi:       data.NoCatatanPolisi.String,
		TanggalCatatanPolisi:     db.Date(data.TanggalCatatanPolisi.Time),
		AkteKelahiran:            data.AkteKelahiran.String,
		NomorBPJS:                data.Bpjs.String,
		NPWP:                     data.Npwp.String,
		TanggalNPWP:              db.Date(data.TanggalNpwp.Time),
		NomorTaspen:              data.NoTaspen.String,
		UnitOrganisasi:           unitOrganisasi,
		Photo:                    data.Foto,
		UnorID:                   data.UnorID,
	}, nil
}

func (s *service) updatePegawai(ctx context.Context, nip string, params updatePegawaiParams) (bool, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repo get pegawai pns id: %w", err)
	}

	if err := s.validatePegawaiReferences(ctx, params); err != nil {
		return false, err
	}

	err = s.withTransaction(ctx, func(txRepo repository) error {
		updateParams := sqlc.UpdateDataPegawaiParams{
			GelarDepan:           params.GelarDepan,
			Nama:                 params.Nama,
			GelarBelakang:        params.GelarBelakang,
			NipBaru:              params.NIP,
			JenisKelamin:         params.JenisKelamin,
			Nik:                  params.NIK,
			Kk:                   params.NomorKK,
			TempatLahirID:        params.TempatLahirID,
			TanggalLahir:         params.TanggalLahir.ToPgtypeDate(),
			TingkatPendidikanID:  typeutil.FromPtr(params.TingkatPendidikanID),
			PendidikanID:         params.PendidikanID,
			JenisKawinID:         typeutil.FromPtr(params.StatusPernikahanID),
			AgamaID:              typeutil.FromPtr(params.AgamaID),
			JenisPegawaiID:       typeutil.FromPtr(params.JenisPegawaiID),
			MasaKerja:            params.MasaKerjaGolongan,
			JenisJabatanID:       typeutil.FromPtr(params.JenisJabatanID),
			JabatanID:            params.JabatanID,
			UnorID:               params.UnitOrganisasiID,
			LokasiKerjaID:        params.LokasiKerjaID,
			GolAwalID:            typeutil.FromPtr(params.GolonganRuangAwalID),
			GolID:                typeutil.FromPtr(params.GolonganRuangAkhirID),
			TmtGolongan:          params.TMTGolongan.ToPgtypeDate(),
			TmtPns:               params.TMTASN.ToPgtypeDate(),
			NoSkCpns:             params.NomorSKASN,
			StatusCpnsPns:        params.StatusPNS,
			EmailDikbud:          params.EmailDikbud,
			Email:                params.EmailPribadi,
			Alamat:               params.Alamat,
			NoHp:                 params.NoHP,
			NoDarurat:            params.NoKontakDarurat,
			NoSuratDokter:        params.NomorSuratDokter,
			TanggalSuratDokter:   params.TanggalSuratDokter.ToPgtypeDate(),
			NoBebasNarkoba:       params.NomorSuratBebasNarkoba,
			TanggalBebasNarkoba:  params.TanggalSuratBebasNarkoba.ToPgtypeDate(),
			NoCatatanPolisi:      params.NomorCatatanPolisi,
			TanggalCatatanPolisi: params.TanggalCatatanPolisi.ToPgtypeDate(),
			AkteKelahiran:        params.AkteKelahiran,
			Bpjs:                 params.NomorBPJS,
			Npwp:                 params.NPWP,
			TanggalNpwp:          params.TanggalNPWP.ToPgtypeDate(),
			NoTaspen:             params.NomorTaspen,
			MkBulan:              typeutil.FromPtr(params.MkBulan),
			MkTahun:              typeutil.FromPtr(params.MkTahun),
			MkBulanSwasta:        typeutil.FromPtr(params.MkBulanSwasta),
			MkTahunSwasta:        typeutil.FromPtr(params.MkTahunSwasta),
			Nip:                  nip,
			PnsID:                pnsID,
		}

		if err := txRepo.UpdateDataPegawai(ctx, updateParams); err != nil {
			return fmt.Errorf("repo update data pegawai: %w", err)
		}

		if err := txRepo.UpdateTTDPegawaiNIPByNIP(ctx, sqlc.UpdateTTDPegawaiNIPByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update ttd pegawai nip: %w", err)
		}

		if err := txRepo.UpdateAnakNIPByNIP(ctx, sqlc.UpdateAnakNIPByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update anak nip: %w", err)
		}

		if err := txRepo.UpdateOrangTuaNIPByNIP(ctx, sqlc.UpdateOrangTuaNIPByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update orang tua nip: %w", err)
		}

		if err := txRepo.UpdatePasanganNIPByNIP(ctx, sqlc.UpdatePasanganNIPByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update pasangan nip: %w", err)
		}

		if err := txRepo.UpdateRiwayatAsesmenNineBoxNamaNipByPNSID(ctx, sqlc.UpdateRiwayatAsesmenNineBoxNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat asesmen nine box: %w", err)
		}

		if err := txRepo.UpdateRiwayatAsesmenNamaNipByPNSID(ctx, sqlc.UpdateRiwayatAsesmenNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat asesmen: %w", err)
		}

		if err := txRepo.UpdateRiwayatHukumanDisiplinNamaNipByPNSID(ctx, sqlc.UpdateRiwayatHukumanDisiplinNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat hukuman disiplin: %w", err)
		}

		if err := txRepo.UpdateRiwayatJabatanNamaNipByPNSID(ctx, sqlc.UpdateRiwayatJabatanNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat jabatan: %w", err)
		}

		if err := txRepo.UpdateRiwayatKenaikanGajiBerkalaNamaNipByPNSID(ctx, sqlc.UpdateRiwayatKenaikanGajiBerkalaNamaNipByPNSIDParams{
			Nama:         params.Nama,
			NipBaru:      params.NIP,
			TanggalLahir: params.TanggalLahir.ToPgtypeDate(),
			Nip:          nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat kenaikan gaji berkala: %w", err)
		}

		if err := txRepo.UpdateRiwayatKepangkatanNamaNipByPNSID(ctx, sqlc.UpdateRiwayatKepangkatanNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat kepangkatan: %w", err)
		}

		if err := txRepo.UpdateRiwayatKinerjaNamaNipByPNSID(ctx, sqlc.UpdateRiwayatKinerjaNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat kinerja: %w", err)
		}

		if err := txRepo.UpdateRiwayatPelatihanFungsionalNamaNipByNIP(ctx, sqlc.UpdateRiwayatPelatihanFungsionalNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat pelatihan fungsional: %w", err)
		}

		if err := txRepo.UpdateRiwayatPelatihanSIASNNamaNipByPNSID(ctx, sqlc.UpdateRiwayatPelatihanSIASNNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat pelatihan siasn: %w", err)
		}

		if err := txRepo.UpdateRiwayatPelatihanStrukturalNamaNipByPNSID(ctx, sqlc.UpdateRiwayatPelatihanStrukturalNamaNipByPNSIDParams{
			NipBaru: params.NIP,
			Nama:    params.Nama,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat pelatihan struktural: %w", err)
		}

		if err := txRepo.UpdateRiwayatPelatihanTeknisNIPByPNSID(ctx, sqlc.UpdateRiwayatPelatihanTeknisNIPByPNSIDParams{
			NipBaru: params.NIP,
			PnsID:   pnsID,
		}); err != nil {
			return fmt.Errorf("repo update riwayat pelatihan teknis: %w", err)
		}

		if err := txRepo.UpdateRiwayatPendidikanNamaNipByNIP(ctx, sqlc.UpdateRiwayatPendidikanNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat pendidikan: %w", err)
		}

		if err := txRepo.UpdateRiwayatPenghargaanNamaNipByNIP(ctx, sqlc.UpdateRiwayatPenghargaanNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat penghargaan: %w", err)
		}

		if err := txRepo.UpdateRiwayatPenugasanNamaNipByNIP(ctx, sqlc.UpdateRiwayatPenugasanNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat penugasan: %w", err)
		}

		if err := txRepo.UpdateRiwayatPindahUnitKerjaNamaNipByNIP(ctx, sqlc.UpdateRiwayatPindahUnitKerjaNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
			Nama:    params.Nama,
		}); err != nil {
			return fmt.Errorf("repo update riwayat pindah unit kerja: %w", err)
		}

		if err := txRepo.UpdateRiwayatSertifikasiNamaNipByNIP(ctx, sqlc.UpdateRiwayatSertifikasiNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat sertifikasi: %w", err)
		}

		if err := txRepo.UpdateRiwayatUjiKompetensiNamaNipByNIP(ctx, sqlc.UpdateRiwayatUjiKompetensiNamaNipByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat uji kompetensi: %w", err)
		}

		if err := txRepo.UpdateRiwayatSuratKeputusanNamaNipPemrosesByNIP(ctx, sqlc.UpdateRiwayatSuratKeputusanNamaNipPemrosesByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update riwayat surat keputusan: %w", err)
		}

		if err := txRepo.UpdateSuratKeputusanNamaNipPemilikByNIP(ctx, sqlc.UpdateSuratKeputusanNamaNipPemilikByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update surat keputusan pemilik: %w", err)
		}

		if err := txRepo.UpdateSuratKeputusanNipPemrosesByNIP(ctx, sqlc.UpdateSuratKeputusanNipPemrosesByNIPParams{
			NipBaru: params.NIP,
			Nip:     nip,
		}); err != nil {
			return fmt.Errorf("repo update surat keputusan pemroses: %w", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("update pegawai transaction: %w", err)
	}

	return true, nil
}

func (s *service) validatePegawaiReferences(ctx context.Context, params updatePegawaiParams) error {
	var errs []error

	if params.AgamaID != nil {
		if _, err := s.repo.GetRefAgama(ctx, int32(*params.AgamaID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get agama: %w", err)
			}
			errs = append(errs, errors.New("data agama tidak ditemukan"))
		}
	}

	if params.GolonganRuangAwalID != nil {
		if _, err := s.repo.GetRefGolongan(ctx, int32(*params.GolonganRuangAwalID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get golongan awal: %w", err)
			}
			errs = append(errs, errors.New("data golongan awal tidak ditemukan"))
		}
	}

	if params.GolonganRuangAkhirID != nil {
		if _, err := s.repo.GetRefGolongan(ctx, int32(*params.GolonganRuangAkhirID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get golongan akhir: %w", err)
			}
			errs = append(errs, errors.New("data golongan akhir tidak ditemukan"))
		}
	}

	if params.StatusPernikahanID != nil {
		if _, err := s.repo.GetRefJenisKawin(ctx, int32(*params.StatusPernikahanID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get status pernikahan: %w", err)
			}
			errs = append(errs, errors.New("data status pernikahan tidak ditemukan"))
		}
	}

	if params.TingkatPendidikanID != nil {
		if _, err := s.repo.GetRefTingkatPendidikan(ctx, int32(*params.TingkatPendidikanID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get tingkat pendidikan: %w", err)
			}
			errs = append(errs, errors.New("data tingkat pendidikan tidak ditemukan"))
		}
	}

	if params.JenisJabatanID != nil {
		if _, err := s.repo.GetRefJenisJabatan(ctx, int32(*params.JenisJabatanID)); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get jenis jabatan: %w", err)
			}
			errs = append(errs, errors.New("data jenis jabatan tidak ditemukan"))
		}
	}

	if len(errs) > 0 {
		return api.NewMultiError(errs)
	}

	return nil
}
