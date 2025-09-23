# kepegawaian.orang_tua

## Description

Orang tua pegawai

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | integer | nextval('orang_tua_id_seq'::regclass) | false |  |  | id data orang tua |
| hubungan | smallint |  | true |  |  | Kode hubungan, 1: ayah, 2: ibu |
| akte_meninggal | varchar(255) |  | true |  |  | Nomor akte meninggal orang tua |
| tanggal_meninggal | date |  | true |  |  | Tanggal meninggal orang tua |
| nama | varchar(255) |  | true |  |  | Nama lengkap orang tua |
| gelar_depan | varchar(20) |  | true |  |  | Gelar di depan nama orang tua |
| gelar_belakang | varchar(50) |  | true |  |  | Gelar di belakang nama orang tua |
| tempat_lahir | varchar(100) |  | true |  |  | Tempat lahir orang tua |
| tanggal_lahir | date |  | true |  |  | Tanggal lahir orang tua |
| agama_id | smallint |  | true |  | [kepegawaian.ref_agama](kepegawaian.ref_agama.md) | id agama orang tua (rujuk ref_agama) |
| email | varchar(255) |  | true |  |  | Alamat email orang tua (bila ada) |
| jenis_dokumen | varchar(10) |  | true |  |  | Jenis dokumen identitas, enum: KTP, PASPOR |
| no_dokumen | varchar(100) |  | true |  |  | Nomor dokumen identitas |
| nip | varchar(20) |  | true |  |  | NIP pegawai |
| pns_id | varchar(36) |  | true |  | [kepegawaian.pegawai](kepegawaian.pegawai.md) | Referensi ke pegawai.pns_id |
| created_at | timestamp with time zone | now() | true |  |  | Waktu perekaman data |
| updated_at | timestamp with time zone | now() | true |  |  | Waktu terakhir pembaruan |
| deleted_at | timestamp with time zone |  | true |  |  | Waktu penghapusan data |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| fk_orang_tua_agama | FOREIGN KEY | FOREIGN KEY (agama_id) REFERENCES ref_agama(id) |
| orang_tua_pkey | PRIMARY KEY | PRIMARY KEY (id) |
| fk_orang_tua_pns_id | FOREIGN KEY | FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| orang_tua_pkey | CREATE UNIQUE INDEX orang_tua_pkey ON kepegawaian.orang_tua USING btree (id) |

## Relations

```mermaid
erDiagram

"kepegawaian.orang_tua" }o--o| "kepegawaian.ref_agama" : "FOREIGN KEY (agama_id) REFERENCES ref_agama(id)"
"kepegawaian.orang_tua" }o--o| "kepegawaian.pegawai" : "FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id)"

"kepegawaian.orang_tua" {
  integer id
  smallint hubungan
  varchar_255_ akte_meninggal
  date tanggal_meninggal
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
"kepegawaian.ref_agama" {
  integer id
  varchar_20_ nama
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
"kepegawaian.pegawai" {
  integer id
  varchar_36_ pns_id
  varchar_9_ nip_lama
  varchar_20_ nip_baru
  varchar_100_ nama
  varchar_50_ gelar_depan
  varchar_50_ gelar_belakang
  varchar_50_ tempat_lahir_id
  date tanggal_lahir
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
  date tanggal_sk_cpns
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
  date tanggal_surat_dokter
  varchar_100_ no_bebas_narkoba
  date tanggal_bebas_narkoba
  varchar_100_ no_catatan_polisi
  date tanggal_catatan_polisi
  varchar_50_ akte_kelahiran
  varchar_15_ status_hidup
  varchar_50_ akte_meninggal
  date tanggal_meninggal
  varchar_100_ no_askes
  varchar_100_ no_taspen
  date tanggal_npwp
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
  smallint status_pegawai_backup
  varchar_50_ masa_kerja
  varchar_50_ kartu_asn
  timestamp_with_time_zone created_at
  timestamp_with_time_zone updated_at
  timestamp_with_time_zone deleted_at
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
