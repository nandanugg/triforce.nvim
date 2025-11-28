BEGIN;

CREATE INDEX IF NOT EXISTS koreksi_surat_keputusan_file_id_idx
    ON koreksi_surat_keputusan (file_id);

CREATE INDEX IF NOT EXISTS koreksi_surat_keputusan_pegawai_korektor_id_idx
    ON koreksi_surat_keputusan (pegawai_korektor_id);

COMMIT;
