BEGIN;

ALTER TABLE ref_template
    ADD COLUMN filename varchar(255) NOT NULL;

COMMIT;
