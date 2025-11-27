BEGIN;

ALTER TABLE riwayat_hukuman_disiplin DROP COLUMN s3_file_id;
ALTER TABLE riwayat_pendidikan DROP COLUMN s3_file_id;
ALTER TABLE riwayat_golongan DROP COLUMN s3_file_id;
ALTER TABLE riwayat_diklat_struktural DROP COLUMN s3_file_id;
ALTER TABLE riwayat_diklat_fungsional DROP COLUMN s3_file_id;
ALTER TABLE riwayat_penugasan DROP COLUMN s3_file_id;
ALTER TABLE riwayat_kenaikan_gaji_berkala DROP COLUMN s3_file_id;
ALTER TABLE riwayat_sertifikasi DROP COLUMN s3_file_id;
ALTER TABLE riwayat_diklat DROP COLUMN s3_file_id;
ALTER TABLE riwayat_kursus DROP COLUMN s3_file_id;
ALTER TABLE riwayat_jabatan DROP COLUMN s3_file_id;
ALTER TABLE riwayat_pindah_unit_kerja DROP COLUMN s3_file_id;
ALTER TABLE surat_keputusan DROP COLUMN s3_file_id;
ALTER TABLE surat_keputusan DROP COLUMN s3_file_sign_id;
ALTER TABLE riwayat_penghargaan_umum DROP COLUMN s3_file_id;

DROP TABLE s3_files;

COMMIT;
