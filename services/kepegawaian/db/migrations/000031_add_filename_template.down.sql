BEGIN;

ALTER TABLE ref_template
    DROP COLUMN filename;

COMMIT;
