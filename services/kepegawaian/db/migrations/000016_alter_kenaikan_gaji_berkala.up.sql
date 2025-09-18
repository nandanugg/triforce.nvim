BEGIN;

ALTER TABLE riwayat_kgb ADD COLUMN gaji_pokok INT;

ALTER TABLE riwayat_kgb RENAME COLUMN n_golongan_id TO golongan_id;
ALTER TABLE riwayat_kgb RENAME COLUMN n_gol_tmt TO tmt_golongan;
ALTER TABLE riwayat_kgb RENAME COLUMN n_masakerja_thn TO masa_kerja_golongan_tahun;
ALTER TABLE riwayat_kgb RENAME COLUMN n_masakerja_bln TO masa_kerja_golongan_bulan;
ALTER TABLE riwayat_kgb RENAME COLUMN n_tmt_jabatan TO tmt_jabatan;
ALTER TABLE riwayat_kgb RENAME COLUMN n_jabatan_text TO jabatan;
ALTER TABLE riwayat_kgb RENAME COLUMN tgl_sk TO tanggal_sk;

ALTER TABLE riwayat_kgb RENAME TO riwayat_kenaikan_gaji_berkala;

COMMIT;