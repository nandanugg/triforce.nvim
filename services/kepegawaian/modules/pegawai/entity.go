package pegawai

import (
	"slices"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type profile struct {
	NIPLama        string      `json:"nip_lama"`
	NIPBaru        string      `json:"nip_baru"`
	GelarDepan     string      `json:"gelar_depan"`
	GelarBelakang  string      `json:"gelar_belakang"`
	Nama           string      `json:"nama"`
	Pangkat        string      `json:"pangkat"`
	Golongan       string      `json:"golongan"`
	Jabatan        string      `json:"jabatan"`
	UnitOrganisasi []string    `json:"unit_organisasi"`
	Photo          pgtype.Text `json:"photo"`
}

type pegawai struct {
	PNSID         string      `json:"pns_id"`
	NIP           string      `json:"nip"`
	GelarDepan    string      `json:"gelar_depan"`
	GelarBelakang string      `json:"gelar_belakang"`
	Nama          string      `json:"nama"`
	Golongan      string      `json:"golongan"`
	Jabatan       string      `json:"jabatan"`
	UnitKerja     string      `json:"unit_kerja"`
	Status        string      `json:"status"`
	Photo         pgtype.Text `json:"photo"`
}

const (
	statusKedudukanHukumMPP = "Masa Persiapan Pensiun"
	statusPNSLabel          = "PNS"
	statusCPNSLabel         = "CPNS"
	tipePegawaiAktif        = "aktif"
	tipePegawaiPPNPN        = "ppnpn"
	tipePegawaiNonAktif     = "nonaktif"
)

var (
	statusPNSInDB  = []string{"PNS", "P"}
	statusCPNSInDB = []string{"CPNS", "C"}
)

func getStatusHukum(params string) string {
	switch params {
	case "MPP":
		return statusKedudukanHukumMPP
	default:
		return ""
	}
}

func getStatusPNSDB(params string) []string {
	switch params {
	case statusPNSLabel:
		return statusPNSInDB
	case statusCPNSLabel:
		return statusCPNSInDB
	default:
		return nil
	}
}

func getLabelStatusPNS(namaKedudukuanHukum, statusPNSCPNS string) string {
	if namaKedudukuanHukum == statusKedudukanHukumMPP {
		return "MPP"
	}
	if slices.Contains(statusPNSInDB, statusPNSCPNS) {
		return statusPNSLabel
	}
	if slices.Contains(statusCPNSInDB, statusPNSCPNS) {
		return statusCPNSLabel
	}

	// for other status, return as-is
	return statusPNSCPNS
}

const (
	jenisJabatanStruktural = 1
)

type pegawaiDetail struct {
	PNSID                    string      `json:"pns_id"`
	Nama                     string      `json:"nama"`
	GelarDepan               string      `json:"gelar_depan"`
	GelarBelakang            string      `json:"gelar_belakang"`
	JabatanAktual            string      `json:"jabatan_aktual"`
	JenisJabatanAktual       string      `json:"jenis_jabatan_aktual"`
	TMTJabatan               db.Date     `json:"tmt_jabatan"`
	NIP                      string      `json:"nip"`
	NIK                      string      `json:"nik"`
	NomorKK                  string      `json:"nomor_kk"`
	JenisKelamin             string      `json:"jenis_kelamin"`
	TempatLahir              string      `json:"tempat_lahir"`
	TanggalLahir             db.Date     `json:"tanggal_lahir"`
	TingkatPendidikan        string      `json:"tingkat_pendidikan"`
	Pendidikan               string      `json:"pendidikan"`
	StatusPerkawinan         string      `json:"status_perkawinan"`
	Agama                    string      `json:"agama"`
	EmailDikbud              string      `json:"email_dikbud"`
	EmailPribadi             string      `json:"email_pribadi"`
	Alamat                   string      `json:"alamat"`
	NomorHP                  string      `json:"nomor_hp"`
	NomorKontakDarurat       string      `json:"nomor_kontak_darurat"`
	JenisPegawai             string      `json:"jenis_pegawai"`
	MasaKerjaKeseluruhan     string      `json:"masa_kerja_keseluruhan"`
	MasaKerjaGolongan        string      `json:"masa_kerja_golongan"`
	Jabatan                  string      `json:"jabatan"`
	JenisJabatan             string      `json:"jenis_jabatan"`
	KelasJabatan             string      `json:"kelas_jabatan"`
	LokasiKerja              string      `json:"lokasi_kerja"`
	GolonganRuangAwal        string      `json:"golongan_ruang_awal"`
	GolonganRuangAkhir       string      `json:"golongan_ruang_akhir"`
	PangkatAkhir             string      `json:"pangkat_akhir"`
	TMTGolongan              db.Date     `json:"tmt_golongan"`
	TMTASN                   db.Date     `json:"tmt_asn"`
	NomorSKASN               string      `json:"nomor_sk_asn"`
	IsPPPK                   bool        `json:"is_pppk"`
	StatusASN                string      `json:"status_asn"`
	StatusPNS                string      `json:"status_pns"`
	TMTPNS                   db.Date     `json:"tmt_pns"`
	KartuPegawai             string      `json:"kartu_pegawai"`
	NomorSuratDokter         string      `json:"nomor_surat_dokter"`
	TanggalSuratDokter       db.Date     `json:"tanggal_surat_dokter"`
	NomorSuratBebasNarkoba   string      `json:"nomor_surat_bebas_narkoba"`
	TanggalSuratBebasNarkoba db.Date     `json:"tanggal_surat_bebas_narkoba"`
	NomorCatatanPolisi       string      `json:"nomor_catatan_polisi"`
	TanggalCatatanPolisi     db.Date     `json:"tanggal_catatan_polisi"`
	AkteKelahiran            string      `json:"akte_kelahiran"`
	NomorBPJS                string      `json:"nomor_bpjs"`
	NPWP                     string      `json:"npwp"`
	TanggalNPWP              db.Date     `json:"tanggal_npwp"`
	NomorTaspen              string      `json:"nomor_taspen"`
	UnitOrganisasi           []string    `json:"unit_organisasi"`
	Photo                    pgtype.Text `json:"photo"`
	UnorID                   pgtype.Text `json:"unor_id"`
}
