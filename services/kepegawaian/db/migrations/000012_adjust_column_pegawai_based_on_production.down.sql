BEGIN;

ALTER TABLE pegawai
  ALTER COLUMN gelar_depan TYPE varchar(20),
  ALTER COLUMN kartu_pegawai TYPE varchar(30),
  ALTER COLUMN nik TYPE varchar(20),
  ALTER COLUMN npwp TYPE varchar(20),
  ALTER COLUMN bpjs TYPE varchar(20),
  ALTER COLUMN alamat TYPE varchar(200),
  ALTER COLUMN jabatan_nama TYPE varchar(200),
  ALTER COLUMN jabatan_instansi_nama TYPE varchar(200);

COMMIT;
