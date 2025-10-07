BEGIN;

ALTER TABLE ref_jabatan
    DROP COLUMN IF EXISTS no;

ALTER TABLE ref_jabatan
    ADD COLUMN IF NOT EXISTS tunjangan_jabatan BIGINT;

CREATE SEQUENCE IF NOT EXISTS ref_jabatan_id_seq OWNED BY ref_jabatan.id;

SELECT setval('ref_jabatan_id_seq', COALESCE(MAX(id), 0) + 1, false)
FROM ref_jabatan;

ALTER TABLE ref_jabatan
    ALTER COLUMN id SET DEFAULT nextval('ref_jabatan_id_seq');

COMMIT;
