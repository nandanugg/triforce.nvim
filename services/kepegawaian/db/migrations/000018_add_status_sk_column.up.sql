BEGIN;

ALTER table file_digital_signature
    ADD COLUMN IF NOT EXISTS status_sk SMALLINT;

COMMIT;