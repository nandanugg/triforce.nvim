BEGIN;

CREATE TABLE ref_agama (
  id serial PRIMARY KEY,
  nama varchar(20),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_golongan (
  id serial PRIMARY KEY,
  nama varchar(10),
  nama_pangkat varchar(50),
  nama_2 varchar(10),
  gol int2,
  gol_pppk varchar(10),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_instansi (
  id varchar(36) PRIMARY KEY,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jabatan (
  kode_jabatan varchar(36) PRIMARY KEY,
  id int4 NOT NULL,
  no int4 NOT NULL,
  nama_jabatan varchar(200),
  nama_jabatan_full varchar(200),
  jenis_jabatan int2,
  kelas int2,
  pensiun int2,
  kode_bkn varchar(36),
  nama_jabatan_bkn varchar(200),
  kategori_jabatan varchar(100),
  bkn_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_diklat (
  id serial PRIMARY KEY,
  bkn_id int2,
  jenis_diklat varchar(50),
  kode varchar(2),
  status int2,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_diklat_fungsional (
  id serial PRIMARY KEY,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_diklat_struktural (
  id serial PRIMARY KEY,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_hukuman (
  id serial PRIMARY KEY,
  dikbud_hr_id VARCHAR(2) GENERATED ALWAYS AS (LPAD(id::TEXT, 2, '0')) STORED,
  nama varchar(100),
  tingkat_hukuman varchar(1),
  nama_tingkat_hukuman varchar(10),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_jabatan (
  id serial PRIMARY KEY,
  nama varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_kawin (
  id serial PRIMARY KEY,
  nama varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_kp (
  id serial PRIMARY KEY,
  dikbud_hr_id varchar(4) GENERATED ALWAYS AS (LPAD(id::TEXT, 4, '0')) STORED ,
  nama varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_pegawai (
  id serial PRIMARY KEY,
  dikbud_hr_id VARCHAR(2) GENERATED ALWAYS AS (LPAD(id::TEXT, 2, '0')) STORED,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_penghargaan (
  id varchar(3) PRIMARY KEY,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_kedudukan_hukum (
  id serial PRIMARY KEY,
  dikbud_hr_id VARCHAR(4) GENERATED ALWAYS AS (LPAD(id::TEXT, 4, '0')) STORED,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_kpkn (
  id varchar(36) PRIMARY KEY,
  nama varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_lokasi (
  id varchar(36) PRIMARY KEY,
  kanreg_id varchar(2),
  lokasi_id varchar(36),
  nama varchar(100),
  jenis varchar(2),
  jenis_kabupaten varchar(3),
  jenis_desa varchar(1),
  ibukota varchar(100),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE ref_jenis_satker (
  id serial PRIMARY KEY,
  nama varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE anak (
  id bigserial PRIMARY KEY,
  pasangan_id int8,
  nama varchar(100),
  jenis_kelamin varchar(1),
  tanggal_lahir date,
  tempat_lahir varchar(100),
  status_anak varchar(1),
  pns_id varchar(36),
  nip varchar(20),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE pasangan (
  id bigserial PRIMARY KEY,
  pns int2,
  nama varchar(100),
  tanggal_menikah date,
  akte_nikah varchar(100),
  tanggal_meninggal date,
  akte_meninggal varchar(100),
  tanggal_cerai date,
  akte_cerai varchar(100),
  karsus varchar(100),
  status int2,
  hubungan int2,
  pns_id varchar(36),
  nip varchar(20),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE orang_tua (
  id serial PRIMARY KEY,
  hubungan int2,
  akte_meninggal varchar(255),
  tgl_meninggal date,
  nama varchar(255),
  gelar_depan varchar(20),
  gelar_belakang varchar(50),
  tempat_lahir varchar(100),
  tanggal_lahir date,
  agama_id int2,
  email varchar(255),
  jenis_dokumen varchar(10),
  no_dokumen varchar(100),
  nip varchar(20),
  pns_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE pegawai (
  id serial PRIMARY KEY,
  pns_id varchar(36) UNIQUE NOT NULL,
  nip_lama varchar(9),
  nip_baru varchar(20),
  nama varchar(100),
  gelar_depan varchar(20),
  gelar_belakang varchar(50),
  tempat_lahir_id varchar(50),
  tgl_lahir date,
  jenis_kelamin varchar(1),
  agama_id int2,
  jenis_kawin_id int2,
  nik varchar(20),
  no_darurat varchar(60),
  no_hp varchar(60),
  email varchar(60),
  alamat varchar(200),
  npwp varchar(20),
  bpjs varchar(20),
  jenis_pegawai_id int2,
  kedudukan_hukum_id int4,
  status_cpns_pns varchar(20),
  kartu_pegawai varchar(30),
  no_sk_cpns varchar(100),
  tgl_sk_cpns date,
  tmt_cpns date,
  tmt_pns date,
  gol_awal_id int2,
  gol_id int2,
  tmt_golongan date,
  mk_tahun int2,
  mk_bulan int2,
  jabatan_id varchar(36),
  tmt_jabatan date,
  pendidikan_id varchar(36),
  tahun_lulus int2,
  kpkn_id varchar(36),
  lokasi_kerja_id varchar(36),
  unor_id varchar(36),
  unor_induk_id varchar(36),
  instansi_induk_id varchar(36),
  instansi_kerja_id varchar(36),
  satuan_kerja_induk_id varchar(36),
  satuan_kerja_kerja_id varchar(36),
  golongan_darah varchar(10),
  foto varchar(200),
  tmt_pensiun date,
  lokasi_kerja varchar(36),
  jml_istri int2,
  jml_anak int2,
  no_surat_dokter varchar(100),
  tgl_surat_dokter date,
  no_bebas_narkoba varchar(100),
  tgl_bebas_narkoba date,
  no_catatan_polisi varchar(100),
  tgl_catatan_polisi date,
  akte_kelahiran varchar(50),
  status_hidup varchar(15),
  akte_meninggal varchar(50),
  tgl_meninggal date,
  no_askes varchar(100),
  no_taspen varchar(100),
  tgl_npwp date,
  tempat_lahir varchar(100),
  tingkat_pendidikan_id int2,
  tempat_lahir_nama varchar(200),
  jenis_jabatan_nama varchar(200),
  jabatan_nama varchar(200),
  kpkn_nama varchar(200),
  instansi_induk_nama varchar(200),
  instansi_kerja_nama varchar(200),
  satuan_kerja_induk_nama varchar(200),
  satuan_kerja_nama varchar(200),
  jabatan_instansi_id int4,
  bup int2 DEFAULT 58,
  jabatan_instansi_nama varchar(200),
  jenis_jabatan_id int2,
  terminated_date date,
  status_pegawai int2 DEFAULT 1,
  jabatan_ppnpn varchar(200),
  jabatan_instansi_real_id int4,
  created_by int4,
  updated_by int4,
  email_dikbud_bak varchar(100),
  email_dikbud varchar(100),
  kodecepat varchar(100),
  is_dosen int2,
  mk_tahun_swasta int2 DEFAULT 0,
  mk_bulan_swasta int2 DEFAULT 0,
  kk varchar(30),
  nidn varchar(30),
  ket varchar(200),
  no_sk_pemberhentian varchar(100),
  status_pegawai_backup int2,
  masa_kerja varchar(50),
  kartu_asn varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE pendidikan (
  id varchar(36) PRIMARY KEY,
  tingkat_pendidikan_id int2,
  nama varchar(200),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE pindah_unit (
  id serial PRIMARY KEY,
  nip varchar(20) NOT NULL,
  surat_permohonan_pindah varchar(200),
  unit_asal varchar(36),
  unit_tujuan varchar(36),
  surat_pernyataan_melepas varchar(200),
  sk_kp_terakhir varchar(100),
  sk_jabatan varchar(100),
  skp varchar(10),
  sk_tunkin varchar(100),
  surat_pernyataan_menerima varchar(200),
  no_sk_pindah varchar(100),
  tanggal_sk_pindah date,
  file_sk varchar(200),
  status_satker int2,
  status_biro int2,
  jabatan_id int2,
  keterangan varchar(200),
  tanggal_tmt_pindah date,
  created_by int4,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_assesmen (
  id serial PRIMARY KEY,
  pns_id varchar(36),
  pns_nip varchar(20),
  tahun int2,
  file_upload varchar(200),
  nilai float4,
  nilai_kinerja float4,
  tahun_penilaian_id int2,
  tahun_penilaian_title varchar(50),
  nama_lengkap varchar(100),
  posisi_id varchar(20),
  unit_org_id varchar(36),
  nama_unor varchar(200),
  saran_pengembangan text,
  file_upload_fb_potensi varchar(200),
  file_upload_lengkap_pt varchar(200),
  file_upload_fb_pt varchar(200),
  file_upload_exists int2 DEFAULT 0,
  satker_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_diklat (
  id bigserial PRIMARY KEY,
  jenis_diklat varchar(200),
  jenis_diklat_id int2,
  institusi_penyelenggara varchar(200),
  no_sertifikat varchar(100),
  tanggal_mulai date,
  tanggal_selesai date,
  tahun_diklat int2,
  durasi_jam int2,
  pns_orang_id varchar(36),
  nip_baru varchar(20),
  diklat_struktural_id varchar(36),
  nama_diklat varchar(200),
  file_base64 text,
  rumpun_diklat_nama varchar(200),
  rumpun_diklat_id varchar(36),
  sudah_kirim_siasn varchar(10) DEFAULT ('belum'::character varying),
  siasn_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_diklat_fungsional (
  id varchar(36) PRIMARY KEY,
  nip_baru varchar(20),
  nip_lama varchar(9),
  jenis_diklat varchar(200),
  nama_kursus varchar(200),
  jumlah_jam int4,
  tahun int2,
  institusi_penyelenggara varchar(200),
  jenis_kursus_sertifikat varchar(1),
  no_sertifikat varchar(100),
  instansi varchar(200),
  status_data varchar(50),
  tanggal_kursus date,
  file_base64 text,
  keterangan_berkas varchar(200),
  lama float4,
  siasn_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_diklat_struktural (
  id varchar(36) PRIMARY KEY,
  pns_id varchar(36),
  pns_nip varchar(20),
  pns_nama varchar(100),
  jenis_diklat_id int4,
  nama_diklat varchar(200),
  nomor varchar(100),
  tanggal date,
  tahun int2,
  status_data varchar(10),
  file_base64 text,
  keterangan_berkas varchar(200),
  lama float4,
  siasn_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_golongan (
  id serial PRIMARY KEY,
  pns_id varchar(36),
  pns_nip varchar(20),
  pns_nama varchar(100),
  kode_jenis_kp varchar(4),
  jenis_kp varchar(50),
  golongan_id int2,
  golongan_nama varchar(10),
  pangkat_nama varchar(50),
  sk_nomor varchar(50),
  no_bkn varchar(100),
  jumlah_angka_kredit_utama int2,
  jumlah_angka_kredit_tambahan int2,
  mk_golongan_tahun int2,
  mk_golongan_bulan int2,
  sk_tanggal date,
  tanggal_bkn date,
  tmt_golongan date,
  status_satker int4,
  status_biro int4,
  pangkat_terakhir int4,
  bkn_id varchar(36),
  file_base64 text,
  keterangan_berkas varchar(200),
  arsip_id int8,
  golongan_asal varchar(2),
  basic varchar(15),
  sk_type int2,
  kanreg varchar(5),
  kpkn varchar(50),
  keterangan varchar(200),
  lpnk varchar(10),
  jenis_riwayat varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_hukdis (
  id bigserial PRIMARY KEY,
  pns_id varchar(36),
  pns_nip varchar(20),
  nama varchar(200),
  golongan_id int2,
  nama_golongan varchar(20),
  jenis_hukuman_id int2,
  nama_jenis_hukuman varchar(100),
  sk_nomor varchar(30),
  sk_tanggal date,
  tanggal_mulai_hukuman date,
  masa_tahun int2,
  masa_bulan int2,
  tanggal_akhir_hukuman date,
  no_pp varchar(100),
  no_sk_pembatalan varchar(100),
  tanggal_sk_pembatalan date,
  bkn_id varchar(255),
  file_base64 text,
  keterangan_berkas varchar(200),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_jabatan (
  bkn_id varchar(36),
  pns_id varchar(36),
  pns_nip varchar(20),
  pns_nama varchar(100),
  unor_id varchar(100),
  unor text,
  jenis_jabatan_id varchar(10),
  jenis_jabatan varchar(250),
  jabatan_id varchar(100),
  nama_jabatan text,
  eselon_id varchar(36),
  eselon varchar(100),
  tmt_jabatan date,
  no_sk varchar(100),
  tanggal_sk date,
  satuan_kerja_id varchar(36),
  tmt_pelantikan date,
  is_active int2,
  eselon1 text,
  eselon2 text,
  eselon3 text,
  eselon4 text,
  id bigserial PRIMARY KEY,
  catatan varchar(200),
  jenis_sk varchar(100),
  status_satker int4,
  status_biro int4,
  jabatan_id_bkn varchar(36),
  unor_id_bkn varchar(36),
  tabel_mutasi_id int8,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_kgb (
  pegawai_id int4,
  tmt_sk date,
  alasan varchar(255),
  mv_kgb_id int8,
  no_sk varchar(100),
  pejabat varchar(255),
  id bigserial PRIMARY KEY,
  ref varchar(255) DEFAULT (public.uuid_generate_v4()),
  tgl_sk date,
  pegawai_nama varchar(255),
  pegawai_nip varchar(20),
  birth_place varchar(255),
  birth_date date,
  n_gol_ruang varchar(50),
  n_gol_tmt date,
  n_masakerja_thn int2,
  n_masakerja_bln int2,
  n_gapok varchar(200),
  n_jabatan_text varchar(200),
  n_tmt_jabatan date,
  n_golongan_id int4,
  unit_kerja_induk_text varchar(200),
  unit_kerja_induk_id varchar(200),
  kantor_pembayaran varchar(200),
  last_education varchar(200),
  last_education_date date,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_kinerja (
  id serial PRIMARY KEY,
  tahun int4,
  nip varchar(20),
  nama varchar(200),
  nip_penilai varchar(20),
  nama_penilai varchar(200),
  jabatan_penilai varchar(200),
  unit_kerja_penilai varchar(200),
  nip_penilai_realisasi varchar(20),
  nama_penilai_realisasi varchar(200),
  jabatan_penilai_realisasi varchar(200),
  unit_kerja_penilai_realisasi varchar(200),
  rating_hasil_kerja varchar(50),
  rating_perilaku_kerja varchar(50),
  predikat_kinerja varchar(100),
  ref uuid DEFAULT (public.uuid_generate_v4()),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_kursus (
  id serial PRIMARY KEY,
  pns_nip varchar(20),
  tipe_kursus varchar(10),
  jenis_kursus varchar(30),
  nama_kursus varchar(200),
  lama_kursus float8,
  tanggal_kursus date,
  no_sertifikat varchar(100),
  instansi varchar(200),
  institusi_penyelenggara varchar(200),
  siasn_id varchar(36),
  pns_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_nine_box (
  id serial PRIMARY KEY,
  pns_nip varchar(20),
  nama varchar(200),
  nama_jabatan varchar(200),
  kelas_jabatan int2,
  kesimpulan varchar(200),
  tahun varchar(4),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_pendidikan (
  id serial PRIMARY KEY,
  pns_id_3 varchar(32),
  tingkat_pendidikan_id int2,
  pendidikan_id_3 varchar(32),
  tanggal_lulus date,
  no_ijazah varchar(100),
  nama_sekolah varchar(200),
  gelar_depan varchar(50),
  gelar_belakang varchar(60),
  pendidikan_pertama varchar(1),
  negara_sekolah varchar(255),
  tahun_lulus varchar(4),
  nip varchar(20),
  diakui_bkn int4,
  status_satker int4,
  status_biro int4,
  pendidikan_terakhir int4,
  pns_id varchar(36),
  pendidikan_id varchar(36),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_penghargaan_umum (
  id serial PRIMARY KEY,
  jenis_penghargaan varchar(50),
  deskripsi_penghargaan varchar(100),
  tanggal_penghargaan date,
  exist bool DEFAULT true,
  file_base64 text,
  nip varchar(20),
  nama_penghargaan varchar(200),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_pindah_unit_kerja (
  id bigserial PRIMARY KEY,
  pns_id varchar(36),
  pns_nip varchar(20),
  pns_nama varchar(100),
  sk_nomor varchar(100),
  asal_id varchar(100),
  asal_nama varchar(100),
  unor_id_baru varchar(36),
  nama_unor_baru varchar(200),
  instansi_id varchar(36),
  nama_instansi varchar(200),
  sk_tanggal date,
  satuan_kerja_id varchar(36),
  nama_satuan_kerja varchar(200),
  file_base64 text,
  keterangan_berkas varchar(200),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_ujikom (
  id bigserial PRIMARY KEY,
  jenis_ujikom varchar(100),
  nip_baru varchar(20),
  link_sertifikat text,
  exist bool,
  tahun int4,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE tingkat_pendidikan (
  id int4 PRIMARY KEY,
  golongan_id int4,
  nama varchar(200),
  golongan_awal_id int4,
  abbreviation varchar(200),
  tingkat int2,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE unit_kerja (
  id varchar(36) PRIMARY KEY,
  no int4,
  kode_internal varchar(36),
  nama_unor varchar(200),
  eselon_id varchar(36),
  cepat_kode varchar(36),
  nama_jabatan varchar(200),
  nama_pejabat varchar(200),
  diatasan_id varchar(36),
  instansi_id varchar(36),
  pemimpin_pns_id varchar(36),
  jenis_unor_id varchar(36),
  unor_induk varchar(36),
  jumlah_ideal_staff int2,
  "order" int4,
  is_satker int2 NOT NULL DEFAULT 0,
  eselon_1 varchar(36),
  eselon_2 varchar(36),
  eselon_3 varchar(36),
  eselon_4 varchar(36),
  expired_date date,
  keterangan varchar(200),
  jenis_satker varchar(200),
  abbreviation varchar(200),
  unor_induk_penyetaraan varchar(200),
  jabatan_id varchar(32),
  waktu varchar(4),
  peraturan varchar(100),
  remark varchar(50),
  aktif bool,
  eselon_nama varchar(50),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE update_mandiri (
  id serial PRIMARY KEY,
  pns_id varchar(36),
  kolom varchar(70),
  dari varchar(400),
  perubahan varchar(400),
  status int4,
  verifikasi_by int4,
  verifikasi_tgl date,
  nama_kolom varchar(100),
  level_update int4,
  tabel_id int4,
  updated_by int4,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE riwayat_sertifikasi (
  id bigserial not null,
  nip varchar(20) null,
  tahun int8 null,
  nama_sertifikasi varchar(255) null,
  file_base64 text null,
  deskripsi text null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

-- Foreign key constraints for the kepegawaian database

-- Table: anak
ALTER TABLE anak
  ADD CONSTRAINT fk_anak_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

-- Table: orang_tua
ALTER TABLE orang_tua
  ADD CONSTRAINT fk_orang_tua_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE orang_tua
  ADD CONSTRAINT fk_orang_tua_agama
  FOREIGN KEY (agama_id) REFERENCES ref_agama(id);

-- Table: pasangan
ALTER TABLE pasangan
  ADD CONSTRAINT fk_pasangan_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

-- Table: pegawai
ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_agama
  FOREIGN KEY (agama_id) REFERENCES ref_agama(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_jenis_kawin
  FOREIGN KEY (jenis_kawin_id) REFERENCES ref_jenis_kawin(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_jabatan
  FOREIGN KEY (jabatan_id) REFERENCES ref_jabatan(kode_jabatan);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_golongan
  FOREIGN KEY (gol_id) REFERENCES ref_golongan(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_golongan_awal
  FOREIGN KEY (gol_awal_id) REFERENCES ref_golongan(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_pendidikan
  FOREIGN KEY (tingkat_pendidikan_id) REFERENCES tingkat_pendidikan(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_kpkn
  FOREIGN KEY (kpkn_id) REFERENCES ref_kpkn(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_lokasi_kerja
  FOREIGN KEY (lokasi_kerja_id) REFERENCES ref_lokasi(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_unor
  FOREIGN KEY (unor_id) REFERENCES unit_kerja(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_instansi_induk
  FOREIGN KEY (instansi_induk_id) REFERENCES ref_instansi(id);

ALTER TABLE pegawai
  ADD CONSTRAINT fk_pegawai_instansi_kerja
  FOREIGN KEY (instansi_kerja_id) REFERENCES ref_instansi(id);

-- Table: pindah_unit
ALTER TABLE pindah_unit
  ADD CONSTRAINT fk_pindah_unit_unit_asal
  FOREIGN KEY (unit_asal) REFERENCES unit_kerja(id);

ALTER TABLE pindah_unit
  ADD CONSTRAINT fk_pindah_unit_unit_tujuan
  FOREIGN KEY (unit_tujuan) REFERENCES unit_kerja(id);

-- Table: pendidikan
ALTER TABLE pendidikan
  ADD CONSTRAINT fk_pendidikan_tingkat
  FOREIGN KEY (tingkat_pendidikan_id) REFERENCES tingkat_pendidikan(id);

-- Table: riwayat_assesmen
ALTER TABLE riwayat_assesmen
  ADD CONSTRAINT fk_riwayat_assesmen_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_assesmen
  ADD CONSTRAINT fk_riwayat_assesmen_unit_org
  FOREIGN KEY (unit_org_id) REFERENCES unit_kerja(id);

-- Table: riwayat_diklat
ALTER TABLE riwayat_diklat
  ADD CONSTRAINT fk_riwayat_diklat_pns_id
  FOREIGN KEY (pns_orang_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_diklat
  ADD CONSTRAINT fk_riwayat_diklat_jenis
  FOREIGN KEY (jenis_diklat_id) REFERENCES ref_jenis_diklat(id);

-- Table: riwayat_diklat_struktural
ALTER TABLE riwayat_diklat_struktural
  ADD CONSTRAINT fk_riwayat_diklat_struktural_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_diklat_struktural
  ADD CONSTRAINT fk_riwayat_diklat_struktural_jenis
  FOREIGN KEY (jenis_diklat_id) REFERENCES ref_jenis_diklat_struktural(id);

-- Table: riwayat_golongan
ALTER TABLE riwayat_golongan
  ADD CONSTRAINT fk_riwayat_golongan_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_golongan
  ADD CONSTRAINT fk_riwayat_golongan_golongan
  FOREIGN KEY (golongan_id) REFERENCES ref_golongan(id);

-- Table: riwayat_hukdis
ALTER TABLE riwayat_hukdis
  ADD CONSTRAINT fk_riwayat_hukdis_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_hukdis
  ADD CONSTRAINT fk_riwayat_hukdis_jenis_hukuman
  FOREIGN KEY (jenis_hukuman_id) REFERENCES ref_jenis_hukuman(id);

-- Table: riwayat_jabatan
ALTER TABLE riwayat_jabatan
  ADD CONSTRAINT fk_riwayat_jabatan_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_jabatan
  ADD CONSTRAINT fk_riwayat_jabatan_satuan_kerja
  FOREIGN KEY (satuan_kerja_id) REFERENCES unit_kerja(id);

-- Table: riwayat_kursus
ALTER TABLE riwayat_kursus
  ADD CONSTRAINT fk_riwayat_kursus_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

-- Table: riwayat_pendidikan
ALTER TABLE riwayat_pendidikan
  ADD CONSTRAINT fk_riwayat_pendidikan_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

ALTER TABLE riwayat_pendidikan
  ADD CONSTRAINT fk_riwayat_pendidikan_tingkat
  FOREIGN KEY (tingkat_pendidikan_id) REFERENCES tingkat_pendidikan(id);

ALTER TABLE riwayat_pendidikan
  ADD CONSTRAINT fk_riwayat_pendidikan_pendidikan
  FOREIGN KEY (pendidikan_id) REFERENCES pendidikan(id);

-- Table: riwayat_pindah_unit_kerja
ALTER TABLE riwayat_pindah_unit_kerja
  ADD CONSTRAINT fk_riwayat_pindah_unit_kerja_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

-- Table: unit_kerja
ALTER TABLE unit_kerja
  ADD CONSTRAINT fk_unit_kerja_diatasan
  FOREIGN KEY (diatasan_id) REFERENCES unit_kerja(id);

ALTER TABLE unit_kerja
  ADD CONSTRAINT fk_unit_kerja_instansi
  FOREIGN KEY (instansi_id) REFERENCES ref_instansi(id);

ALTER TABLE unit_kerja
  ADD CONSTRAINT fk_unit_kerja_pemimpin
  FOREIGN KEY (pemimpin_pns_id) REFERENCES pegawai(pns_id);

-- Table: update_mandiri
ALTER TABLE update_mandiri
  ADD CONSTRAINT fk_update_mandiri_pns_id
  FOREIGN KEY (pns_id) REFERENCES pegawai(pns_id);

COMMIT;
