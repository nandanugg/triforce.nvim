BEGIN;

ALTER TABLE ref_tingkat_pendidikan 
  ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE IF EXISTS ref_tingkat_pendidikan_id_seq;

COMMIT;