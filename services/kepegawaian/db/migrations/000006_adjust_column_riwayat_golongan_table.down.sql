BEGIN;

ALTER TABLE riwayat_golongan
DROP CONSTRAINT IF EXISTS fk_riwayat_golongan_jenis_kp;

ALTER TABLE riwayat_golongan
DROP COLUMN IF EXISTS jenis_kp_id;

COMMIT;