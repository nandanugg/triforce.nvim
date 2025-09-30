BEGIN;

ALTER TABLE ref_kedudukan_hukum
    ADD COLUMN is_pegawai_aktif BOOLEAN NOT NULL DEFAULT TRUE;

COMMIT;