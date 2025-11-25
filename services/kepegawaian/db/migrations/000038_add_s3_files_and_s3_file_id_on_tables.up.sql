BEGIN;

CREATE TABLE s3_files (
    id bigserial PRIMARY KEY,
    object_key varchar(255),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    deleted_at timestamptz
);

ALTER TABLE riwayat_hukuman_disiplin ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_pendidikan ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_golongan ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_diklat_struktural ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_diklat_fungsional ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_penugasan ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_kenaikan_gaji_berkala ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_sertifikasi ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_diklat ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_kursus ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_jabatan ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_pindah_unit_kerja ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE surat_keputusan ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);
ALTER TABLE surat_keputusan ADD COLUMN s3_file_sign_id bigint REFERENCES s3_files(id);
ALTER TABLE riwayat_penghargaan_umum ADD COLUMN s3_file_id bigint REFERENCES s3_files(id);

COMMIT;
