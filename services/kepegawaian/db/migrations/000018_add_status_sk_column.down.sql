BEGIN;

ALTER table file_digital_signature
    DROP COLUMN IF EXISTS status_sk;

COMMIT;