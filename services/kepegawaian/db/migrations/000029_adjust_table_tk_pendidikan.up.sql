BEGIN;

CREATE SEQUENCE IF NOT EXISTS ref_tingkat_pendidikan_id_seq OWNED BY ref_tingkat_pendidikan.id;

ALTER TABLE ref_tingkat_pendidikan 
  ALTER COLUMN id SET DEFAULT nextval('ref_tingkat_pendidikan_id_seq');

SELECT setval('ref_tingkat_pendidikan_id_seq', COALESCE(MAX(id), 0) + 1, false)
FROM ref_tingkat_pendidikan;

COMMIT;