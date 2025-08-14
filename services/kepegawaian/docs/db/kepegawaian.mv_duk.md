# kepegawaian.mv_duk

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE MATERIALIZED VIEW mv_duk AS (
 SELECT vw."NAMA_UNOR",
    pegawai."JENIS_JABATAN_ID",
    pegawai."JABATAN_ID",
    jabatan."NAMA_JABATAN",
    pegawai."NIP_LAMA",
    pegawai."NIP_BARU",
    pegawai."NAMA",
    pegawai."GELAR_DEPAN",
    pegawai."GELAR_BELAKANG",
    vw."ESELON_ID" AS vw_eselon_id,
    pegawai."GOL_ID",
    (((golongan."NAMA_PANGKAT")::text || ' '::text) || (golongan."NAMA")::text) AS golongan_text,
    'jabatanku'::text AS jabatan_text,
    pegawai."PNS_ID",
    (((date_part('year'::text, (now())::date) - date_part('year'::text, pegawai."TGL_LAHIR")) * (12)::double precision) + (date_part('month'::text, (now())::date) - date_part('month'::text, pegawai."TGL_LAHIR"))) AS bulan_usia,
    '#'::text AS separator,
    pegawai."TEMPAT_LAHIR_ID",
    pegawai."TGL_LAHIR",
    pegawai."JENIS_KELAMIN",
    pegawai."AGAMA_ID",
    pegawai."JENIS_KAWIN_ID",
    pegawai."NIK",
    pegawai."NOMOR_DARURAT",
    pegawai."NOMOR_HP",
    pegawai."EMAIL",
    pegawai."ALAMAT",
    pegawai."NPWP",
    pegawai."BPJS",
    pegawai."JENIS_PEGAWAI_ID",
    pegawai."KEDUDUKAN_HUKUM_ID",
    pegawai."STATUS_CPNS_PNS",
    pegawai."KARTU_PEGAWAI",
    pegawai."NOMOR_SK_CPNS",
    pegawai."TGL_SK_CPNS",
    pegawai."TMT_CPNS",
    pegawai."TMT_PNS",
    pegawai."GOL_AWAL_ID",
    pegawai."TMT_GOLONGAN",
    pegawai."MK_TAHUN",
    pegawai."MK_BULAN",
    pegawai."TMT_JABATAN",
    pegawai."PENDIDIKAN_ID",
    pegawai."PENDIDIKAN",
    pegawai."TAHUN_LULUS",
    pegawai."KPKN_ID",
    pegawai."LOKASI_KERJA_ID",
    pegawai."UNOR_ID",
    pegawai."UNOR_INDUK_ID",
    pegawai."INSTANSI_INDUK_ID",
    pegawai."INSTANSI_KERJA_ID",
    pegawai."SATUAN_KERJA_INDUK_ID",
    pegawai."SATUAN_KERJA_KERJA_ID",
    pegawai."GOLONGAN_DARAH",
    pegawai."ID",
    pegawai."PHOTO",
    pegawai."TMT_PENSIUN",
    pegawai."BUP",
    vw."ESELON_1",
    vw."ESELON_2",
    vw."ESELON_3",
    vw."ESELON_4",
    vw."NAMA_UNOR_ESELON_4",
    vw."NAMA_UNOR_ESELON_3",
    vw."NAMA_UNOR_ESELON_2",
    vw."NAMA_UNOR_ESELON_1"
   FROM ((((kepegawaian.pns_aktif pa
     LEFT JOIN kepegawaian.pegawai pegawai ON ((pa."ID" = pegawai."ID")))
     LEFT JOIN kepegawaian.golongan ON (((pegawai."GOL_ID")::text = (golongan."ID")::text)))
     LEFT JOIN kepegawaian.vw_unit_list vw ON (((vw."ID")::text = (pegawai."UNOR_ID")::text)))
     LEFT JOIN kepegawaian.jabatan jabatan ON ((pegawai."JABATAN_INSTANSI_ID" = (jabatan."KODE_JABATAN")::bpchar)))
  WHERE ((pegawai."KEDUDUKAN_HUKUM_ID")::text <> ALL (ARRAY[('14'::character varying)::text, ('52'::character varying)::text, ('66'::character varying)::text, ('67'::character varying)::text, ('77'::character varying)::text, ('78'::character varying)::text, ('98'::character varying)::text, ('99'::character varying)::text]))
  ORDER BY pegawai."JENIS_JABATAN_ID", vw."ESELON_ID", vw."ESELON_1", vw."ESELON_2", vw."ESELON_3", vw."ESELON_4", pegawai."JABATAN_ID", vw."NAMA_UNOR_FULL", pegawai."GOL_ID" DESC, pegawai."TMT_GOLONGAN", pegawai."TMT_JABATAN", pegawai."TMT_CPNS", pegawai."TGL_LAHIR"
)
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| NAMA_UNOR | varchar(255) |  | true |  |  |  |
| JENIS_JABATAN_ID | integer |  | true |  |  |  |
| JABATAN_ID | varchar(36) |  | true |  |  |  |
| NAMA_JABATAN | varchar |  | true |  |  |  |
| NIP_LAMA | varchar(9) |  | true |  |  |  |
| NIP_BARU | varchar(18) |  | true |  |  |  |
| NAMA | varchar(255) |  | true |  |  |  |
| GELAR_DEPAN | varchar(60) |  | true |  |  |  |
| GELAR_BELAKANG | varchar(60) |  | true |  |  |  |
| vw_eselon_id | varchar(255) |  | true |  |  |  |
| GOL_ID | integer |  | true |  |  |  |
| golongan_text | text |  | true |  |  |  |
| jabatan_text | text |  | true |  |  |  |
| PNS_ID | varchar(36) |  | true |  |  |  |
| bulan_usia | double precision |  | true |  |  |  |
| separator | text |  | true |  |  |  |
| TEMPAT_LAHIR_ID | varchar(50) |  | true |  |  |  |
| TGL_LAHIR | date |  | true |  |  |  |
| JENIS_KELAMIN | varchar(10) |  | true |  |  |  |
| AGAMA_ID | integer |  | true |  |  |  |
| JENIS_KAWIN_ID | varchar(255) |  | true |  |  |  |
| NIK | varchar(255) |  | true |  |  |  |
| NOMOR_DARURAT | varchar(255) |  | true |  |  |  |
| NOMOR_HP | varchar(60) |  | true |  |  |  |
| EMAIL | varchar(255) |  | true |  |  |  |
| ALAMAT | varchar(255) |  | true |  |  |  |
| NPWP | varchar(255) |  | true |  |  |  |
| BPJS | varchar(50) |  | true |  |  |  |
| JENIS_PEGAWAI_ID | varchar(50) |  | true |  |  |  |
| KEDUDUKAN_HUKUM_ID | varchar(36) |  | true |  |  |  |
| STATUS_CPNS_PNS | varchar(20) |  | true |  |  |  |
| KARTU_PEGAWAI | varchar(30) |  | true |  |  |  |
| NOMOR_SK_CPNS | varchar(60) |  | true |  |  |  |
| TGL_SK_CPNS | date |  | true |  |  |  |
| TMT_CPNS | date |  | true |  |  |  |
| TMT_PNS | date |  | true |  |  |  |
| GOL_AWAL_ID | varchar(36) |  | true |  |  |  |
| TMT_GOLONGAN | date |  | true |  |  |  |
| MK_TAHUN | varchar(20) |  | true |  |  |  |
| MK_BULAN | varchar(20) |  | true |  |  |  |
| TMT_JABATAN | date |  | true |  |  |  |
| PENDIDIKAN_ID | varchar(36) |  | true |  |  |  |
| PENDIDIKAN | character(165) |  | true |  |  |  |
| TAHUN_LULUS | varchar(20) |  | true |  |  |  |
| KPKN_ID | varchar(36) |  | true |  |  |  |
| LOKASI_KERJA_ID | varchar(36) |  | true |  |  |  |
| UNOR_ID | varchar(36) |  | true |  |  |  |
| UNOR_INDUK_ID | varchar(36) |  | true |  |  |  |
| INSTANSI_INDUK_ID | varchar(36) |  | true |  |  |  |
| INSTANSI_KERJA_ID | varchar(36) |  | true |  |  |  |
| SATUAN_KERJA_INDUK_ID | varchar(36) |  | true |  |  |  |
| SATUAN_KERJA_KERJA_ID | varchar(36) |  | true |  |  |  |
| GOLONGAN_DARAH | varchar(20) |  | true |  |  |  |
| ID | integer |  | true |  |  |  |
| PHOTO | varchar(100) |  | true |  |  |  |
| TMT_PENSIUN | date |  | true |  |  |  |
| BUP | smallint |  | true |  |  |  |
| ESELON_1 | varchar(32) |  | true |  |  |  |
| ESELON_2 | varchar(32) |  | true |  |  |  |
| ESELON_3 | varchar(32) |  | true |  |  |  |
| ESELON_4 | varchar(32) |  | true |  |  |  |
| NAMA_UNOR_ESELON_4 | varchar(255) |  | true |  |  |  |
| NAMA_UNOR_ESELON_3 | varchar(255) |  | true |  |  |  |
| NAMA_UNOR_ESELON_2 | varchar(255) |  | true |  |  |  |
| NAMA_UNOR_ESELON_1 | varchar(255) |  | true |  |  |  |

## Referenced Tables

| Name | Columns | Comment | Type |
| ---- | ------- | ------- | ---- |
| [kepegawaian.pns_aktif](kepegawaian.pns_aktif.md) | 3 |  | VIEW |
| [kepegawaian.pegawai](kepegawaian.pegawai.md) | 100 |  | BASE TABLE |
| [kepegawaian.golongan](kepegawaian.golongan.md) | 6 |  | BASE TABLE |
| [kepegawaian.vw_unit_list](kepegawaian.vw_unit_list.md) | 30 |  | MATERIALIZED VIEW |
| [kepegawaian.jabatan](kepegawaian.jabatan.md) | 12 |  | BASE TABLE |

## Relations

```mermaid
erDiagram


"kepegawaian.mv_duk" {
  varchar_255_ NAMA_UNOR
  integer JENIS_JABATAN_ID
  varchar_36_ JABATAN_ID
  varchar NAMA_JABATAN
  varchar_9_ NIP_LAMA
  varchar_18_ NIP_BARU
  varchar_255_ NAMA
  varchar_60_ GELAR_DEPAN
  varchar_60_ GELAR_BELAKANG
  varchar_255_ vw_eselon_id
  integer GOL_ID
  text golongan_text
  text jabatan_text
  varchar_36_ PNS_ID
  double_precision bulan_usia
  text separator
  varchar_50_ TEMPAT_LAHIR_ID
  date TGL_LAHIR
  varchar_10_ JENIS_KELAMIN
  integer AGAMA_ID
  varchar_255_ JENIS_KAWIN_ID
  varchar_255_ NIK
  varchar_255_ NOMOR_DARURAT
  varchar_60_ NOMOR_HP
  varchar_255_ EMAIL
  varchar_255_ ALAMAT
  varchar_255_ NPWP
  varchar_50_ BPJS
  varchar_50_ JENIS_PEGAWAI_ID
  varchar_36_ KEDUDUKAN_HUKUM_ID
  varchar_20_ STATUS_CPNS_PNS
  varchar_30_ KARTU_PEGAWAI
  varchar_60_ NOMOR_SK_CPNS
  date TGL_SK_CPNS
  date TMT_CPNS
  date TMT_PNS
  varchar_36_ GOL_AWAL_ID
  date TMT_GOLONGAN
  varchar_20_ MK_TAHUN
  varchar_20_ MK_BULAN
  date TMT_JABATAN
  varchar_36_ PENDIDIKAN_ID
  character_165_ PENDIDIKAN
  varchar_20_ TAHUN_LULUS
  varchar_36_ KPKN_ID
  varchar_36_ LOKASI_KERJA_ID
  varchar_36_ UNOR_ID
  varchar_36_ UNOR_INDUK_ID
  varchar_36_ INSTANSI_INDUK_ID
  varchar_36_ INSTANSI_KERJA_ID
  varchar_36_ SATUAN_KERJA_INDUK_ID
  varchar_36_ SATUAN_KERJA_KERJA_ID
  varchar_20_ GOLONGAN_DARAH
  integer ID
  varchar_100_ PHOTO
  date TMT_PENSIUN
  smallint BUP
  varchar_32_ ESELON_1
  varchar_32_ ESELON_2
  varchar_32_ ESELON_3
  varchar_32_ ESELON_4
  varchar_255_ NAMA_UNOR_ESELON_4
  varchar_255_ NAMA_UNOR_ESELON_3
  varchar_255_ NAMA_UNOR_ESELON_2
  varchar_255_ NAMA_UNOR_ESELON_1
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
