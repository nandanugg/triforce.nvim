BEGIN;

ALTER TABLE pegawai
  ALTER COLUMN gelar_depan TYPE varchar(50),
  ALTER COLUMN kartu_pegawai TYPE varchar(50),
  ALTER COLUMN nik TYPE varchar(50),
  ALTER COLUMN npwp TYPE varchar(50),
  ALTER COLUMN bpjs TYPE varchar(50),
  ALTER COLUMN alamat TYPE varchar(300),
  ALTER COLUMN jabatan_nama TYPE varchar(300),
  ALTER COLUMN jabatan_instansi_nama TYPE varchar(400);

COMMIT;
