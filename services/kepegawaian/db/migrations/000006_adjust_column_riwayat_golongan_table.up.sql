BEGIN;

ALTER TABLE riwayat_golongan
ADD COLUMN jenis_kp_id int;

ALTER TABLE riwayat_golongan
ADD CONSTRAINT fk_riwayat_golongan_jenis_kp
FOREIGN KEY (jenis_kp_id) REFERENCES ref_jenis_kp(id);

COMMIT;