package datapribadi

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) getDataPribadi(ctx context.Context, userID int64) (*dataPribadi, error) {
	q := `
		select
			p."ID",
			p."NIP_LAMA",
			p."NIP_BARU",
			p."NIK",
			p."TEMPAT_LAHIR",
			p."TGL_LAHIR",
			p."EMAIL_DIKBUD",
			p."EMAIL",
			p."ALAMAT",
			p."NOMOR_HP",
			a."NAMA",
			p."JENIS_KELAMIN",
			p."TK_PENDIDIKAN",
			p."PENDIDIKAN",
			p."MASA_KERJA",
			p."LOKASI_KERJA",
			jp."NAMA",
			'TODO: PangkatGolonganAktif',
			'TODO: GolonganRuangAwal',
			'TODO: GolonganRuangTerakhir',
			p."TMT_GOLONGAN",
			p."JABATAN_NAMA",
			uk."NAMA_UNOR",
			'TODO: GajiPokok',
			'1990-01-01'::date,
			'TODO: SKASN',
			'TODO: StatusPNS',
			p."KARTU_PEGAWAI",
			jk."NAMA",
			p."NO_SURAT_DOKTER",
			p."TGL_SURAT_DOKTER",
			p."NO_BEBAS_NARKOBA",
			p."TGL_BEBAS_NARKOBA",
			p."NO_CATATAN_POLISI",
			p."TGL_CATATAN_POLISI",
			p."AKTE_KELAHIRAN",
			p."BPJS",
			p."NPWP",
			p."TGL_NPWP",
			p."NOMOR_DARURAT",
			p."NAMA",
			p."PHOTO",
			p."GELAR_DEPAN",
			p."GELAR_BELAKANG"
		from pegawai p
		left join agama a on p."AGAMA_ID" = a."ID"
		left join users u on p."NIP_BARU" = u.nip
		left join jenis_pegawai jp on p."JENIS_PEGAWAI_ID" = jp."ID"
		left join jenis_kawin jk on p."JENIS_KAWIN_ID" = jk."ID"
		left join unitkerja uk on p."UNOR_ID" = uk."ID"
		where u.id = $1
	`

	var data dataPribadi
	err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&data.ID,
		&data.NIP,
		&data.NIPBaru,
		&data.NIK,
		&data.TempatLahir,
		&data.TanggalLahir,
		&data.EmailDikbud,
		&data.EmailLain,
		&data.Alamat,
		&data.NomorHP,
		&data.Agama,
		&data.JenisKelamin,
		&data.TingkatPendidikan,
		&data.Pendidikan,
		&data.MasaKerja,
		&data.LokasiKerja,
		&data.JenisPegawai,
		&data.PangkatGolonganAktif,
		&data.GolonganRuangAwal,
		&data.GolonganRuangTerakhir,
		&data.TMTGolongan,
		&data.Jabatan,
		&data.UnitKerja,
		&data.GajiPokok,
		&data.TMTASN,
		&data.SKASN,
		&data.StatusPNS,
		&data.KartuPegawai,
		&data.StatusPerkawinan,
		&data.NomorSuratDokter,
		&data.TanggalSuratDokter,
		&data.NomorSuratBebasNarkoba,
		&data.TanggalSuratBebasNarkoba,
		&data.NomorCatatanPolisi,
		&data.TanggalCatatanPolisi,
		&data.AkteKelahiran,
		&data.NomorBPJS,
		&data.NPWP,
		&data.TanggalNPWP,
		&data.NomorDarurat,
		&data.Nama,
		&data.Photo,
		&data.GelarDepan,
		&data.GelarBelakang,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	for _, toTrim := range []*string{
		&data.AkteKelahiran,
		&data.LokasiKerja,
		&data.NomorCatatanPolisi,
		&data.NomorSuratBebasNarkoba,
		&data.NomorSuratDokter,
		&data.Pendidikan,
		&data.TempatLahir,
		&data.Jabatan,
	} {
		*toTrim = strings.TrimSpace(*toTrim)
	}

	return &data, err
}

func (r *repository) listStatusPernikahan(ctx context.Context) ([]statusPernikahan, error) {
	rows, err := r.db.QueryContext(ctx, `select "ID", "NAMA" from jenis_kawin order by 2 asc`)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []statusPernikahan{}
	for rows.Next() {
		var row statusPernikahan
		err := rows.Scan(&row.ID, &row.Nama)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}
