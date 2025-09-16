BEGIN;

ALTER TABLE riwayat_penghargaan_umum
    DROP COLUMN jenis_penghargaan;

ALTER TABLE riwayat_penghargaan_umum
    ADD COLUMN jenis_penghargaan_id int REFERENCES ref_jenis_penghargaan(id);

COMMIT;