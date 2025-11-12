BEGIN;

ALTER TABLE riwayat_sertifikasi
DROP CONSTRAINT IF EXISTS riwayat_sertifikasi_pkey;

COMMIT;