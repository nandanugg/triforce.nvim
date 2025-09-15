BEGIN;

ALTER TABLE ref_jenis_penghargaan 
  DROP CONSTRAINT ref_jenis_penghargaan_pkey;

ALTER TABLE ref_jenis_penghargaan 
  ALTER COLUMN id TYPE integer USING id::integer;

CREATE SEQUENCE ref_jenis_penghargaan_id_seq OWNED BY ref_jenis_penghargaan.id;

SELECT setval('ref_jenis_penghargaan_id_seq', COALESCE(MAX(id), 0) + 1, false)
FROM ref_jenis_penghargaan;

ALTER TABLE ref_jenis_penghargaan 
  ALTER COLUMN id SET DEFAULT nextval('ref_jenis_penghargaan_id_seq');

ALTER TABLE ref_jenis_penghargaan 
  ADD PRIMARY KEY (id);

COMMIT;
