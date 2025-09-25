BEGIN;

ALTER TABLE file_digital_signature_riwayat
    RENAME COLUMN pemroses_id to nip_pemroses;

COMMIT;