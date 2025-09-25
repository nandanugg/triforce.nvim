BEGIN;

ALTER TABLE file_digital_signature_riwayat
    RENAME COLUMN nip_pemroses TO pemroses_id;

COMMIT;