BEGIN;

ALTER TABLE ref_jenis_penghargaan 
  DROP CONSTRAINT ref_jenis_penghargaan_pkey;

ALTER TABLE ref_jenis_penghargaan 
  ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE IF EXISTS ref_jenis_penghargaan_id_seq;

ALTER TABLE ref_jenis_penghargaan 
  ALTER COLUMN id TYPE varchar(3);

ALTER TABLE ref_jenis_penghargaan 
  ADD PRIMARY KEY (id);

COMMIT;
