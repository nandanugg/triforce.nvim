# kepegawaian.pegawai

## Description

Data utama pegawai

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | integer | nextval('pegawai_id_seq'::regclass) | false |  |  | Identitas numerik unik baris data pegawai |
| pns_id | varchar(36) |  | false | [kepegawaian.anak](kepegawaian.anak.md) [kepegawaian.pasangan](kepegawaian.pasangan.md) [kepegawaian.orang_tua](kepegawaian.orang_tua.md) [kepegawaian.riwayat_assesmen](kepegawaian.riwayat_assesmen.md) [kepegawaian.riwayat_diklat](kepegawaian.riwayat_diklat.md) [kepegawaian.riwayat_diklat_struktural](kepegawaian.riwayat_diklat_struktural.md) [kepegawaian.riwayat_golongan](kepegawaian.riwayat_golongan.md) [kepegawaian.riwayat_hukdis](kepegawaian.riwayat_hukdis.md) [kepegawaian.riwayat_jabatan](kepegawaian.riwayat_jabatan.md) [kepegawaian.riwayat_kursus](kepegawaian.riwayat_kursus.md) [kepegawaian.riwayat_pendidikan](kepegawaian.riwayat_pendidikan.md) [kepegawaian.riwayat_pindah_unit_kerja](kepegawaian.riwayat_pindah_unit_kerja.md) [kepegawaian.unit_kerja](kepegawaian.unit_kerja.md) [kepegawaian.update_mandiri](kepegawaian.update_mandiri.md) |  | id pegawai negeri sipil (UUID) yang menjadi kunci rujukan antar tabel |
| nip_lama | varchar(9) |  | true |  |  | Nomor Induk Pegawai format lama |
| nip_baru | varchar(20) |  | true |  |  | Nomor Induk Pegawai format baru (20 digit) |
| nama | varchar(100) |  | true |  |  | Nama lengkap pegawai |
| gelar_depan | varchar(50) |  | true |  |  | Gelar akademik/jabatan di depan nama |
| gelar_belakang | varchar(50) |  | true |  |  | Gelar akademik/jabatan di belakang nama |
| tempat_lahir_id | varchar(50) |  | true |  |  | id tempat lahir (rujuk ref_lokasi) |
| tgl_lahir | date |  | true |  |  | Tanggal lahir pegawai |
| jenis_kelamin | varchar(1) |  | true |  |  | Kode jenis kelamin, M: laki-laki, F: perempuan |
| agama_id | smallint |  | true |  | [kepegawaian.ref_agama](kepegawaian.ref_agama.md) | Kode agama (rujuk ref_agama) |
| jenis_kawin_id | smallint |  | true |  | [kepegawaian.ref_jenis_kawin](kepegawaian.ref_jenis_kawin.md) | Status perkawinan (rujuk ref_jenis_kawin) |
| nik | varchar(50) |  | true |  |  | Nomor Induk Kependudukan |
| no_darurat | varchar(60) |  | true |  |  | Nomor telepon yang dapat dihubungi dalam keadaan darurat |
| no_hp | varchar(60) |  | true |  |  | Nomor telepon seluler pegawai |
| email | varchar(60) |  | true |  |  | Alamat surat elektronik pribadi |
| alamat | varchar(300) |  | true |  |  | Alamat domisili pegawai |
| npwp | varchar(50) |  | true |  |  | Nomor Pokok Wajib Pajak |
| bpjs | varchar(50) |  | true |  |  | Nomor kepesertaan BPJS |
| jenis_pegawai_id | smallint |  | true |  |  | id jenis pegawai (PNS/PPPK/dll; rujuk ref_jenis_pegawai) |
| kedudukan_hukum_id | integer |  | true |  |  | id kedudukan hukum (rujuk ref_kedudukan_hukum) |
| status_cpns_pns | varchar(20) |  | true |  |  | Status CPNS/PNS |
| kartu_pegawai | varchar(50) |  | true |  |  | Nomor kartu pegawai |
| no_sk_cpns | varchar(100) |  | true |  |  | Nomor SK pengangkatan CPNS |
| tgl_sk_cpns | date |  | true |  |  | Tanggal SK pengangkatan CPNS |
| tmt_cpns | date |  | true |  |  | Tanggal mulai tugas (CPNS) |
| tmt_pns | date |  | true |  |  | Tanggal mulai tugas (PNS) |
| gol_awal_id | smallint |  | true |  | [kepegawaian.ref_golongan](kepegawaian.ref_golongan.md) | Golongan awal saat pengangkatan (rujuk ref_golongan) |
| gol_id | smallint |  | true |  | [kepegawaian.ref_golongan](kepegawaian.ref_golongan.md) | Golongan terakhir/aktif (rujuk ref_golongan) |
| tmt_golongan | date |  | true |  |  | Tanggal mulai berlaku golongan saat ini |
| mk_tahun | smallint |  | true |  |  | Masa kerja tahun |
| mk_bulan | smallint |  | true |  |  | Masa kerja bulan |
| jabatan_id | varchar(36) |  | true |  | [kepegawaian.ref_jabatan](kepegawaian.ref_jabatan.md) | id jabatan pegawai (rujuk ref_jabatan) |
| tmt_jabatan | date |  | true |  |  | Tanggal mulai jabatan |
| pendidikan_id | varchar(36) |  | true |  |  | id pendidikan (rujuk ref_pendidikan) |
| tahun_lulus | smallint |  | true |  |  | Tahun kelulusan pendidikan terakhir |
| kpkn_id | varchar(36) |  | true |  | [kepegawaian.ref_kpkn](kepegawaian.ref_kpkn.md) | id KPPN/KPKN pembayaran gaji (rujuk ref_kpkn) |
| lokasi_kerja_id | varchar(36) |  | true |  | [kepegawaian.ref_lokasi](kepegawaian.ref_lokasi.md) | id lokasi kerja (rujuk ref_lokasi) |
| unor_id | varchar(36) |  | true |  | [kepegawaian.unit_kerja](kepegawaian.unit_kerja.md) | Unit organisasi/kerja (rujuk unit_kerja) |
| unor_induk_id | varchar(36) |  | true |  |  | Unit organisasi/kerja induk (rujuk unit_kerja) |
| instansi_induk_id | varchar(36) |  | true |  | [kepegawaian.ref_instansi](kepegawaian.ref_instansi.md) | id instansi induk pegawai (rujuk ref_instansi) |
| instansi_kerja_id | varchar(36) |  | true |  | [kepegawaian.ref_instansi](kepegawaian.ref_instansi.md) | id instansi tempat bekerja (rujuk ref_instansi) |
| satuan_kerja_induk_id | varchar(36) |  | true |  |  | id satuan kerja induk pegawai |
| satuan_kerja_kerja_id | varchar(36) |  | true |  |  | id satuan kerja pegawai |
| golongan_darah | varchar(10) |  | true |  |  | Golongan darah |
| foto | varchar(200) |  | true |  |  | Lokasi/URL berkas foto pegawai |
| tmt_pensiun | date |  | true |  |  | Tanggal perkiraan/penetapan pensiun (BUP) |
| lokasi_kerja | varchar(36) |  | true |  |  | Nama lokasi kerja |
| jml_istri | smallint |  | true |  |  | Jumlah pasangan |
| jml_anak | smallint |  | true |  |  | Jumlah anak yang tercatat |
| no_surat_dokter | varchar(100) |  | true |  |  | Nomor surat pemeriksaan kesehatan |
| tgl_surat_dokter | date |  | true |  |  | Tanggal surat pemeriksaan kesehatan |
| no_bebas_narkoba | varchar(100) |  | true |  |  | Nomor Surat Keterangan Bebas Narkoba |
| tgl_bebas_narkoba | date |  | true |  |  | Tanggal Surat Keterangan Bebas Narkoba |
| no_catatan_polisi | varchar(100) |  | true |  |  | Nomor Surat Catatan Kelakukan Baik dari kepolisian |
| tgl_catatan_polisi | date |  | true |  |  | Tanggal Surat Catatan Kelakukan Baik dari kepolisian |
| akte_kelahiran | varchar(50) |  | true |  |  | Nomor akte kelahiran |
| status_hidup | varchar(15) |  | true |  |  | Status hidup pegawai |
| akte_meninggal | varchar(50) |  | true |  |  | Nomor akte meninggal |
| tgl_meninggal | date |  | true |  |  | Tanggal meninggal pegawai |
| no_askes | varchar(100) |  | true |  |  | Nomor ASKES (jika tersedia/legacy) |
| no_taspen | varchar(100) |  | true |  |  | Nomor Taspen |
| tgl_npwp | date |  | true |  |  | Tanggal terbit NPWP |
| tempat_lahir | varchar(100) |  | true |  |  | Nama tempat lahir berdasarkan referensi ref_lokasi |
| tingkat_pendidikan_id | smallint |  | true |  | [kepegawaian.ref_tingkat_pendidikan](kepegawaian.ref_tingkat_pendidikan.md) | Tingkat pendidikan terakhir (rujuk tingkat_pendidikan) |
| tempat_lahir_nama | varchar(200) |  | true |  |  | Nama tempat lahir (teks bebas) |
| jenis_jabatan_nama | varchar(200) |  | true |  |  | Nama jenis jabatan |
| jabatan_nama | varchar(300) |  | true |  |  | Nama jabatan pegawai |
| kpkn_nama | varchar(200) |  | true |  |  | Nama KPPN/KPKN pembayaran gaji |
| instansi_induk_nama | varchar(200) |  | true |  |  | Nama instansi induk pegawai |
| instansi_kerja_nama | varchar(200) |  | true |  |  | Nama instansi tempat bekerja pegawai |
| satuan_kerja_induk_nama | varchar(200) |  | true |  |  | Nama satuan kerja induk pegawai |
| satuan_kerja_nama | varchar(200) |  | true |  |  | Nama satuan kerja pegawai |
| jabatan_instansi_id | varchar(36) |  | true |  | [kepegawaian.ref_jabatan](kepegawaian.ref_jabatan.md) | id jabatan instansi pegawai (rujuk ref_jabatan) |
| bup | smallint | 58 | true |  |  | Batas usia pensiun |
| jabatan_instansi_nama | varchar(400) |  | true |  |  | Nama jabatan instansi pegawai |
| jenis_jabatan_id | smallint |  | true |  |  | id jenis jabatan |
| terminated_date | date |  | true |  |  |  |
| status_pegawai | smallint | 1 | true |  |  | Status pegawai, 1: pns, 2: honorer |
| jabatan_ppnpn | varchar(200) |  | true |  |  | Nama jabatan Pegawai Pemerintah Non Pegawai Negeri |
| jabatan_instansi_real_id | varchar(36) |  | true |  | [kepegawaian.ref_jabatan](kepegawaian.ref_jabatan.md) | id jabatan instansi pegawai (rujuk ref_jabatan) |
| created_by | integer |  | true |  |  | id user yang memasukkan data pegawai |
| updated_by | integer |  | true |  |  | id user yang memperbarui data pegawai |
| email_dikbud_bak | varchar(100) |  | true |  |  | Alamat surat elektronik backup untuk kepentingan pekerjaan |
| email_dikbud | varchar(100) |  | true |  |  | Alamat surat elektronik untuk kepentingan pekerjaan |
| kodecepat | varchar(100) |  | true |  |  |  |
| is_dosen | smallint |  | true |  |  |  |
| mk_tahun_swasta | smallint | 0 | true |  |  | Masa kerja tahun di swasta, sebelum menjadi ASN |
| mk_bulan_swasta | smallint | 0 | true |  |  | Masa kerja bulan di swasta, sebelum menjadi ASN |
| kk | varchar(30) |  | true |  |  | Nomor kartu keluarga |
| nidn | varchar(30) |  | true |  |  |  |
| ket | varchar(200) |  | true |  |  | Keterangan tambahan terhadap pegawai |
| no_sk_pemberhentian | varchar(100) |  | true |  |  | Nomor SK pemberhentian dari PNS |
| status_pegawai_backup | integer |  | true |  |  | Status pegawai backup |
| masa_kerja | varchar(50) |  | true |  |  | masa kerja |
| kartu_asn | varchar(50) |  | true |  |  | Nomor kartu ASN |
| created_at | timestamp with time zone | now() | true |  |  | Waktu perekaman data dibuat |
| updated_at | timestamp with time zone | now() | true |  |  | Waktu terakhir data diperbarui |
| deleted_at | timestamp with time zone |  | true |  |  | Waktu penghapusan lunak (soft delete) bila ada |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| pegawai_id_not_null | n | NOT NULL id |
| pegawai_pns_id_not_null | n | NOT NULL pns_id |
| fk_pegawai_agama | FOREIGN KEY | FOREIGN KEY (agama_id) REFERENCES ref_agama(id) |
| fk_pegawai_golongan | FOREIGN KEY | FOREIGN KEY (gol_id) REFERENCES ref_golongan(id) |
| fk_pegawai_golongan_awal | FOREIGN KEY | FOREIGN KEY (gol_awal_id) REFERENCES ref_golongan(id) |
| fk_pegawai_instansi_induk | FOREIGN KEY | FOREIGN KEY (instansi_induk_id) REFERENCES ref_instansi(id) |
| fk_pegawai_instansi_kerja | FOREIGN KEY | FOREIGN KEY (instansi_kerja_id) REFERENCES ref_instansi(id) |
| fk_pegawai_jabatan | FOREIGN KEY | FOREIGN KEY (jabatan_id) REFERENCES ref_jabatan(kode_jabatan) |
| fk_pegawai_jabatan_instansi | FOREIGN KEY | FOREIGN KEY (jabatan_instansi_id) REFERENCES ref_jabatan(kode_jabatan) |
| fk_pegawai_jabatan_instansi_real | FOREIGN KEY | FOREIGN KEY (jabatan_instansi_real_id) REFERENCES ref_jabatan(kode_jabatan) |
| fk_pegawai_jenis_kawin | FOREIGN KEY | FOREIGN KEY (jenis_kawin_id) REFERENCES ref_jenis_kawin(id) |
| fk_pegawai_kpkn | FOREIGN KEY | FOREIGN KEY (kpkn_id) REFERENCES ref_kpkn(id) |
| fk_pegawai_lokasi_kerja | FOREIGN KEY | FOREIGN KEY (lokasi_kerja_id) REFERENCES ref_lokasi(id) |
| pegawai_pkey | PRIMARY KEY | PRIMARY KEY (id) |
| pegawai_pns_id_key | UNIQUE | UNIQUE (pns_id) |
| fk_pegawai_pendidikan | FOREIGN KEY | FOREIGN KEY (tingkat_pendidikan_id) REFERENCES ref_tingkat_pendidikan(id) |
| fk_pegawai_unor | FOREIGN KEY | FOREIGN KEY (unor_id) REFERENCES unit_kerja(id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| pegawai_pkey | CREATE UNIQUE INDEX pegawai_pkey ON kepegawaian.pegawai USING btree (id) |
| pegawai_pns_id_key | CREATE UNIQUE INDEX pegawai_pns_id_key ON kepegawaian.pegawai USING btree (pns_id) |

## Relations

```mermaid
erDiagram

"kepegawaian.anak" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.pasangan" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.orang_tua" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_assesmen" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_diklat" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_orang_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_diklat_struktural" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_golongan" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_hukdis" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_jabatan" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_kursus" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_pendidikan" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.riwayat_pindah_unit_kerja" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.unit_kerja" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pemimpin_pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.update_mandiri" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_agama" : "FOREIGN KEY (agama_id) REFERENCES ref_agama(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_jenis_kawin" : "FOREIGN KEY (jenis_kawin_id) REFERENCES ref_jenis_kawin(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_golongan" : "FOREIGN KEY (gol_awal_id) REFERENCES ref_golongan(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_golongan" : "FOREIGN KEY (gol_id) REFERENCES ref_golongan(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_jabatan" : "FOREIGN KEY (jabatan_id) REFERENCES ref_jabatan(kode_jabatan)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_kpkn" : "FOREIGN KEY (kpkn_id) REFERENCES ref_kpkn(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_lokasi" : "FOREIGN KEY (lokasi_kerja_id) REFERENCES ref_lokasi(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.unit_kerja" : "FOREIGN KEY (unor_id) REFERENCES unit_kerja(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_instansi" : "FOREIGN KEY (instansi_induk_id) REFERENCES ref_instansi(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_instansi" : "FOREIGN KEY (instansi_kerja_id) REFERENCES ref_instansi(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_tingkat_pendidikan" : "FOREIGN KEY (tingkat_pendidikan_id) REFERENCES ref_tingkat_pendidikan(id)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_jabatan" : "FOREIGN KEY (jabatan_instansi_id) REFERENCES ref_jabatan(kode_jabatan)"
"kepegawaian.pegawai" }o--o| "kepegawaian.ref_jabatan" : "FOREIGN KEY (jabatan_instansi_real_id) REFERENCES ref_jabatan(kode_jabatan)"

"kepegawaian.pegawai" {
  integer id
  varchar_36_ pns_id
  varchar_9_ nip_lama
  varchar_20_ nip_baru
  varchar_100_ nama
  varchar_50_ gelar_depan
  varchar_50_ gelar_belakang
  varchar_50_ tempat_lahir_id
  date tgl_lahir
  varchar_1_ jenis_kelamin
  smallint agama_id FK
  smallint jenis_kawin_id FK
  varchar_50_ nik
  varchar_60_ no_darurat
  varchar_60_ no_hp
  varchar_60_ email
  varchar_300_ alamat
  varchar_50_ npwp
  varchar_50_ bpjs
  smallint jenis_pegawai_id
  integer kedudukan_hukum_id
  varchar_20_ status_cpns_pns
  varchar_50_ kartu_pegawai
  varchar_100_ no_sk_cpns
  date tgl_sk_cpns
  date tmt_cpns
  date tmt_pns
  smallint gol_awal_id FK
  smallint gol_id FK
  date tmt_golongan
  smallint mk_tahun
  smallint mk_bulan
  varchar_36_ jabatan_id FK
  date tmt_jabatan
  varchar_36_ pendidikan_id
  smallint tahun_lulus
  varchar_36_ kpkn_id FK
  varchar_36_ lokasi_kerja_id FK
  varchar_36_ unor_id FK
  varchar_36_ unor_induk_id
  varchar_36_ instansi_induk_id FK
  varchar_36_ instansi_kerja_id FK
  varchar_36_ satuan_kerja_induk_id
  varchar_36_ satuan_kerja_kerja_id
  varchar_10_ golongan_darah
  varchar_200_ foto
  date tmt_pensiun
  varchar_36_ lokasi_kerja
  smallint jml_istri
  smallint jml_anak
  varchar_100_ no_surat_dokter
  date tgl_surat_dokter
  varchar_100_ no_bebas_narkoba
  date tgl_bebas_narkoba
  varchar_100_ no_catatan_polisi
  date tgl_catatan_polisi
  varchar_50_ akte_kelahiran
  varchar_15_ status_hidup
  varchar_50_ akte_meninggal
  date tgl_meninggal
  varchar_100_ no_askes
  varchar_100_ no_taspen
  date tgl_npwp
  varchar_100_ tempat_lahir
  smallint tingkat_pendidikan_id FK
  varchar_200_ tempat_lahir_nama
  varchar_200_ jenis_jabatan_nama
  varchar_300_ jabatan_nama
  varchar_200_ kpkn_nama
  varchar_200_ instansi_induk_nama
  varchar_200_ instansi_kerja_nama
  varchar_200_ satuan_kerja_induk_nama
  varchar_200_ satuan_kerja_nama
  varchar_36_ jabatan_instansi_id FK
  smallint bup
  varchar_400_ jabatan_instansi_nama
  smallint jenis_jabatan_id
  date terminated_date
  smallint status_pegawai
  varchar_200_ jabatan_ppnpn
  varchar_36_ jabatan_instansi_real_id FK
  integer created_by
  integer updated_by
  varchar_100_ email_dikbud_bak
  varchar_100_ email_dikbud
  varchar_100_ kodecepat
  smallint is_dosen
  smallint mk_tahun_swasta
  smallint mk_bulan_swasta
  varchar_30_ kk
  varchar_30_ nidn
  varchar_200_ ket
  varchar_100_ no_sk_pemberhentian
  integer status_pegawai_backup
  varchar_50_ masa_kerja
  varchar_50_ kartu_asn
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.anak" {
  bigint id
  bigint pasangan_id
  varchar_100_ nama
  varchar_1_ jenis_kelamin
  date tanggal_lahir
  varchar_100_ tempat_lahir
  varchar_1_ status_anak
  varchar_36_ pns_id FK
  varchar_20_ nip
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.pasangan" {
  bigint id
  smallint pns
  varchar_100_ nama
  date tanggal_menikah
  varchar_100_ akte_nikah
  date tanggal_meninggal
  varchar_100_ akte_meninggal
  date tanggal_cerai
  varchar_100_ akte_cerai
  varchar_100_ karsus
  smallint status
  smallint hubungan
  varchar_36_ pns_id FK
  varchar_20_ nip
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.orang_tua" {
  integer id
  smallint hubungan
  varchar_255_ akte_meninggal
  date tgl_meninggal
  varchar_255_ nama
  varchar_20_ gelar_depan
  varchar_50_ gelar_belakang
  varchar_100_ tempat_lahir
  date tanggal_lahir
  smallint agama_id FK
  varchar_255_ email
  varchar_10_ jenis_dokumen
  varchar_100_ no_dokumen
  varchar_20_ nip
  varchar_36_ pns_id FK
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.riwayat_assesmen" {
  integer id
  varchar_36_ pns_id FK
  varchar_20_ pns_nip
  smallint tahun
  varchar_200_ file_upload
  real nilai
  real nilai_kinerja
  smallint tahun_penilaian_id
  varchar_50_ tahun_penilaian_title
  varchar_100_ nama_lengkap
  varchar_20_ posisi_id
  varchar_36_ unit_org_id FK
  varchar_200_ nama_unor
  text saran_pengembangan
  varchar_200_ file_upload_fb_potensi
  varchar_200_ file_upload_lengkap_pt
  varchar_200_ file_upload_fb_pt
  smallint file_upload_exists
  varchar_36_ satker_id
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.riwayat_diklat" {
  bigint id
  varchar_200_ jenis_diklat
  smallint jenis_diklat_id FK
  varchar_200_ institusi_penyelenggara
  varchar_100_ no_sertifikat
  date tanggal_mulai
  date tanggal_selesai
  smallint tahun_diklat
  smallint durasi_jam
  varchar_36_ pns_orang_id FK
  varchar_20_ nip_baru
  varchar_36_ diklat_struktural_id
  varchar_200_ nama_diklat
  text file_base64
  varchar_200_ rumpun_diklat_nama
  varchar_36_ rumpun_diklat_id
  varchar_10_ sudah_kirim_siasn
  varchar_36_ siasn_id
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.riwayat_diklat_struktural" {
  varchar_36_ id
  varchar_36_ pns_id FK
  varchar_20_ pns_nip
  varchar_100_ pns_nama
  integer jenis_diklat_id FK
  varchar_200_ nama_diklat
  varchar_100_ nomor
  date tanggal
  smallint tahun
  varchar_10_ status_data
  text file_base64
  varchar_200_ keterangan_berkas
  real lama
  varchar_36_ siasn_id
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.riwayat_golongan" {
  integer id
  varchar_36_ pns_id FK
  varchar_20_ pns_nip
  varchar_100_ pns_nama
  varchar_4_ kode_jenis_kp
  varchar_50_ jenis_kp
  smallint golongan_id FK
  varchar_10_ golongan_nama
  varchar_50_ pangkat_nama
  varchar_50_ sk_nomor
  varchar_100_ no_bkn
  smallint jumlah_angka_kredit_utama
  smallint jumlah_angka_kredit_tambahan
  smallint mk_golongan_tahun
  smallint mk_golongan_bulan
  date sk_tanggal
  date tanggal_bkn
  date tmt_golongan
  integer status_satker
  integer status_biro
  integer pangkat_terakhir
  varchar_36_ bkn_id
  text file_base64
  varchar_200_ keterangan_berkas
  bigint arsip_id
  varchar_2_ golongan_asal
  varchar_15_ basic
  smallint sk_type
  varchar_5_ kanreg
  varchar_50_ kpkn
  varchar_200_ keterangan
  varchar_10_ lpnk
  varchar_50_ jenis_riwayat
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
  integer jenis_kp_id FK
}
"kepegawaian.riwayat_hukdis" {
  bigint id
  varchar_36_ pns_id FK
  varchar_20_ pns_nip
  varchar_200_ nama
  smallint golongan_id
  varchar_20_ nama_golongan
  smallint jenis_hukuman_id FK
  varchar_100_ nama_jenis_hukuman
  varchar_30_ sk_nomor
  date sk_tanggal
  date tanggal_mulai_hukuman
  smallint masa_tahun
  smallint masa_bulan
  date tanggal_akhir_hukuman
  varchar_100_ no_pp
  varchar_100_ no_sk_pembatalan
  date tanggal_sk_pembatalan
  varchar_255_ bkn_id
  text file_base64
  varchar_200_ keterangan_berkas
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.riwayat_jabatan" {
  varchar_36_ bkn_id
  varchar_36_ pns_id FK
  varchar_20_ pns_nip
  varchar_100_ pns_nama
  varchar_36_ unor_id
  text unor
  integer jenis_jabatan_id
  varchar_250_ jenis_jabatan
  integer jabatan_id
  text nama_jabatan
  varchar_36_ eselon_id
  varchar_100_ eselon
  date tmt_jabatan
  varchar_100_ no_sk
  date tanggal_sk
  varchar_36_ satuan_kerja_id FK
  date tmt_pelantikan
  smallint is_active
  text eselon1
  text eselon2
  text eselon3
  text eselon4
  bigint id
  varchar_200_ catatan
  varchar_100_ jenis_sk
  integer status_satker
  integer status_biro
  varchar_36_ jabatan_id_bkn
  varchar_36_ unor_id_bkn
  bigint tabel_mutasi_id
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
  boolean status_plt
  integer kelas_jabatan_id FK
  date periode_jabatan_start_date
  date periode_jabatan_end_date
}
"kepegawaian.riwayat_kursus" {
  integer id
  varchar_20_ pns_nip
  varchar_10_ tipe_kursus
  varchar_30_ jenis_kursus
  varchar_200_ nama_kursus
  double_precision lama_kursus
  date tanggal_kursus
  varchar_100_ no_sertifikat
  varchar_200_ instansi
  varchar_200_ institusi_penyelenggara
  varchar_36_ siasn_id
  varchar_36_ pns_id FK
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.riwayat_pendidikan" {
  integer id
  varchar_32_ pns_id_3
  smallint tingkat_pendidikan_id FK
  varchar_32_ pendidikan_id_3
  date tanggal_lulus
  varchar_100_ no_ijazah
  varchar_200_ nama_sekolah
  varchar_50_ gelar_depan
  varchar_60_ gelar_belakang
  varchar_1_ pendidikan_pertama
  varchar_255_ negara_sekolah
  smallint tahun_lulus
  varchar_20_ nip
  integer diakui_bkn
  integer status_satker
  integer status_biro
  integer pendidikan_terakhir
  varchar_36_ pns_id FK
  varchar_36_ pendidikan_id FK
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
  smallint tugas_belajar
  text file_base64
  varchar_200_ keterangan_berkas
}
"kepegawaian.riwayat_pindah_unit_kerja" {
  bigint id
  varchar_36_ pns_id FK
  varchar_20_ pns_nip
  varchar_100_ pns_nama
  varchar_100_ sk_nomor
  varchar_100_ asal_id
  varchar_100_ asal_nama
  varchar_36_ unor_id_baru
  varchar_200_ nama_unor_baru
  varchar_36_ instansi_id
  varchar_200_ nama_instansi
  date sk_tanggal
  varchar_36_ satuan_kerja_id
  varchar_200_ nama_satuan_kerja
  text file_base64
  varchar_200_ keterangan_berkas
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.unit_kerja" {
  varchar_36_ id
  integer no
  varchar_36_ kode_internal
  varchar_200_ nama_unor
  varchar_36_ eselon_id
  varchar_36_ cepat_kode
  varchar_200_ nama_jabatan
  varchar_200_ nama_pejabat
  varchar_36_ diatasan_id FK
  varchar_36_ instansi_id FK
  varchar_36_ pemimpin_pns_id FK
  varchar_36_ jenis_unor_id
  varchar_36_ unor_induk
  smallint jumlah_ideal_staff
  integer order
  smallint is_satker
  varchar_36_ eselon_1
  varchar_36_ eselon_2
  varchar_36_ eselon_3
  varchar_36_ eselon_4
  date expired_date
  varchar_200_ keterangan
  varchar_200_ jenis_satker
  varchar_200_ abbreviation
  varchar_200_ unor_induk_penyetaraan
  varchar_32_ jabatan_id
  varchar_4_ waktu
  varchar_100_ peraturan
  varchar_50_ remark
  boolean aktif
  varchar_50_ eselon_nama
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.update_mandiri" {
  integer id
  varchar_36_ pns_id FK
  varchar_70_ kolom
  varchar_400_ dari
  varchar_400_ perubahan
  integer status
  integer verifikasi_by
  date verifikasi_tgl
  varchar_100_ nama_kolom
  integer level_update
  integer tabel_id
  integer updated_by
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_agama" {
  integer id
  varchar_20_ nama
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_jenis_kawin" {
  integer id
  varchar_50_ nama
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_golongan" {
  integer id
  varchar_10_ nama
  varchar_50_ nama_pangkat
  varchar_10_ nama_2
  smallint gol
  varchar_10_ gol_pppk
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_jabatan" {
  varchar_36_ kode_jabatan
  integer id
  integer no
  varchar_200_ nama_jabatan
  varchar_200_ nama_jabatan_full
  smallint jenis_jabatan
  smallint kelas
  smallint pensiun
  varchar_36_ kode_bkn
  varchar_200_ nama_jabatan_bkn
  varchar_100_ kategori_jabatan
  varchar_36_ bkn_id
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_kpkn" {
  varchar_36_ id
  varchar_100_ nama
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_lokasi" {
  varchar_36_ id
  varchar_2_ kanreg_id
  varchar_36_ lokasi_id
  varchar_100_ nama
  varchar_2_ jenis
  varchar_3_ jenis_kabupaten
  varchar_1_ jenis_desa
  varchar_100_ ibukota
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_instansi" {
  varchar_36_ id
  varchar_100_ nama
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.ref_tingkat_pendidikan" {
  integer id
  integer golongan_id
  varchar_200_ nama
  integer golongan_awal_id
  varchar_200_ abbreviation
  smallint tingkat
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
