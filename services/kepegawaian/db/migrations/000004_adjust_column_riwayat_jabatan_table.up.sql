BEGIN;

ALTER TABLE riwayat_jabatan
  ALTER COLUMN unor_id TYPE varchar(36),
  ALTER COLUMN jabatan_id TYPE int USING jabatan_id::int,
  ALTER COLUMN jenis_jabatan_id TYPE int USING jenis_jabatan_id::int;
  

ALTER TABLE riwayat_jabatan
  ADD COLUMN status_plt boolean,
  ADD COLUMN kelas_jabatan_id int REFERENCES ref_kelas_jabatan(id),
  ADD COLUMN periode_jabatan_start_date date,
  ADD COLUMN periode_jabatan_end_date date;

COMMIT;