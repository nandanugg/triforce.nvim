BEGIN;

alter table surat_keputusan 
    ALTER COLUMN halaman_ttd TYPE boolean using (CASE WHEN halaman_ttd = 1 THEN true ELSE false END);

DROP INDEX IF EXISTS pegawai_ttd_nip_idx;
DROP TABLE IF EXISTS pegawai_ttd;

COMMIT;