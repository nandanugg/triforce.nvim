BEGIN;

ALTER TABLE riwayat_diklat_struktural
  ADD COLUMN institusi_penyelenggara varchar(200);

COMMIT;
