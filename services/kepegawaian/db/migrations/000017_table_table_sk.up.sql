BEGIN;

CREATE TABLE IF NOT EXISTS file_digital_signature (
  file_id VARCHAR(200) PRIMARY KEY NOT NULL,
  kategori VARCHAR(100),
  file_base64 TEXT,
  ttd_pegawai_id VARCHAR(255),
  status_ttd SMALLINT,
  nip_sk VARCHAR(50),
  no_sk VARCHAR(50),
  tanggal_sk DATE,
  tmt_sk DATE,
  lokasi_file TEXT,
  status_koreksi SMALLINT,
  catatan TEXT,
  pegawai_korektor_id VARCHAR(100),
  asal_surat_sk VARCHAR(100),
  status_kembali SMALLINT,
  nama_pemilik_sk VARCHAR(200),
  jabatan_pemilik_sk TEXT,
  file_base64_sign TEXT,
  unit_kerja_pemilik_sk TEXT,
  nip_pemroses VARCHAR(50),
  ds_ok boolean,
  arsip VARCHAR(50),
  status_pns VARCHAR(20),
  tmt_sampai_dengan DATE, -- khusus untuk Surat Perintah PLT/PLH
  telah_kirim boolean, -- Jika 1, tampilkan di dikbudHR
  halaman_ttd boolean default true, -- halaman diletakan tandataangan digital
  show_qrcode boolean default false, -- 0/null : tidak tampilkan (seperti semula), 1 : tampilkan qrdari bssn
  letak_ttd SMALLINT DEFAULT 0, -- 1:tengah bawah, 2 : kiri Bawah 0: kanan bawah
  kode_unit_kerja_internal VARCHAR(200), -- untuk menampung nama unit kerja internal via kode
  kode_jabatan_internal VARCHAR(200), -- untuk menampung nama jabatan dengan kode jabatan internal
  kelompok_jabatan VARCHAR(200), -- khusus untuk keperluan laporan rekap
  tanggal_ttd timestamptz, -- untuk mengetahui tgl tandatangan
  email_kirim VARCHAR(200), -- Untuk menentukan alamat alternatif pengiriman dokumen
  sent_to_siasin VARCHAR(100) DEFAULT 'n',
  blockchain_issuer_id TEXT,
  blockchain_image_url TEXT,
  blockchain_hash TEXT,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

COMMENT ON COLUMN file_digital_signature.status_koreksi IS '1 sudah dikoreksi, 0 belum dikoreksi, 3 dikembalikan, 2 siap dikoreksi';
COMMENT ON COLUMN file_digital_signature.tmt_sampai_dengan IS 'khusus untuk Surat Perintah PLT/PLH';
COMMENT ON COLUMN file_digital_signature.telah_kirim IS 'Jika 1, tampilkan di dikbudHR';
COMMENT ON COLUMN file_digital_signature.halaman_ttd IS 'halaman diletakan tandataangan digital';
COMMENT ON COLUMN file_digital_signature.show_qrcode IS '0/null : tidak tampilkan (seperti semula), 1 : tampilkan qrdari bssn';
COMMENT ON COLUMN file_digital_signature.letak_ttd IS '1:tengah bawah, 2 : kiri Bawah 0: kanan bawah';
COMMENT ON COLUMN file_digital_signature.kode_unit_kerja_internal IS 'untuk menampung nama unit kerja internal via kode';
COMMENT ON COLUMN file_digital_signature.kode_jabatan_internal IS 'untuk menampung nama jabatan dengan kode jabatan internal';
COMMENT ON COLUMN file_digital_signature.kelompok_jabatan IS 'khusus untuk keperluan laporan rekap';
COMMENT ON COLUMN file_digital_signature.tanggal_ttd IS 'untuk mengetahui tgl tandatangan';
COMMENT ON COLUMN file_digital_signature.email_kirim IS 'Untuk menentukan alamat alternatif pengiriman dokumen';
COMMENT ON COLUMN file_digital_signature.ds_ok IS '1 tanda tangan elektronik, 0 tanda tangan manual';

CREATE TABLE IF NOT EXISTS file_digital_signature_corrector (
  id SERIAL PRIMARY KEY,
  korektor_ke SMALLINT,
  pegawai_korektor_id VARCHAR(100),
  status_kembali SMALLINT, -- 1=dikembalikan, 0/null = sudah oke
  catatan_koreksi TEXT,
  status_koreksi SMALLINT, -- 1=koreksi ok, 2=siap koreksi, 0/null = masih antrian
  file_id VARCHAR(200),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

COMMENT ON COLUMN file_digital_signature_corrector.status_kembali IS '1=dikembalikan, 0/null = sudah oke';
COMMENT ON COLUMN file_digital_signature_corrector.status_koreksi IS '1=koreksi ok, 2=siap koreksi, 0/null = masih antrian';

CREATE TABLE IF NOT EXISTS file_digital_signature_riwayat (
  id BIGSERIAL PRIMARY KEY,
  file_id VARCHAR(200),
  pemroses_id VARCHAR(255),
  tindakan TEXT,
  catatan_tindakan TEXT,
  waktu_tindakan timestamptz,
  akses_pengguna VARCHAR(200),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS log_digital_signature (
  id SERIAL PRIMARY KEY,
  file_id VARCHAR(32),
  nik VARCHAR(30),
  keterangan VARCHAR(255),
  status SMALLINT, -- 1:gagal, 2:berhasil
  proses_cron boolean default false, -- 0 = belum, 1 = sudah
  created_by INTEGER,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

COMMENT ON COLUMN log_digital_signature.status IS '1:gagal, 2:berhasil';
COMMENT ON COLUMN log_digital_signature.proses_cron IS '0 = belum, 1 = sudah';


COMMIT;