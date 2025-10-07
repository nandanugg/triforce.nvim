BEGIN;

ALTER TABLE ref_jabatan
    ADD COLUMN IF NOT EXISTS no integer;

ALTER TABLE ref_jabatan
    DROP COLUMN IF EXISTS tunjangan_jabatan;

ALTER TABLE ref_jabatan
    ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE IF EXISTS ref_jabatan_id_seq;

COMMIT;
