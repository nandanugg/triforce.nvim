BEGIN;

ALTER TABLE riwayat_jabatan
  ALTER COLUMN unor_id TYPE varchar(100);

ALTER TABLE riwayat_jabatan
  DROP COLUMN status_plt,
  DROP COLUMN kelas_jabatan_id,
  DROP COLUMN periode_jabatan_start_date,
  DROP COLUMN periode_jabatan_end_date;

COMMIT;