BEGIN;

ALTER TABLE ref_kedudukan_hukum
    DROP COLUMN is_pegawai_aktif;

COMMIT;