BEGIN;

ALTER TABLE riwayat_penghargaan_umum
    DROP COLUMN jenis_penghargaan_id;

ALTER TABLE riwayat_penghargaan_umum
    ADD COLUMN jenis_penghargaan varchar(255);

COMMIT;