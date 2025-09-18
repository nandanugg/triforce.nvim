BEGIN;

ALTER TABLE riwayat_kenaikan_gaji_berkala RENAME TO riwayat_kgb;

ALTER TABLE riwayat_kgb RENAME COLUMN golongan_id TO n_golongan_id;
ALTER TABLE riwayat_kgb RENAME COLUMN tmt_golongan TO n_gol_tmt;
ALTER TABLE riwayat_kgb RENAME COLUMN masa_kerja_golongan_tahun TO n_masakerja_thn;
ALTER TABLE riwayat_kgb RENAME COLUMN masa_kerja_golongan_bulan TO n_masakerja_bln;
ALTER TABLE riwayat_kgb RENAME COLUMN tmt_jabatan TO n_tmt_jabatan;
ALTER TABLE riwayat_kgb RENAME COLUMN jabatan TO n_jabatan_text;
ALTER TABLE riwayat_kgb RENAME COLUMN tanggal_sk TO tgl_sk;

ALTER TABLE riwayat_kgb DROP COLUMN gaji_pokok;

COMMIT;