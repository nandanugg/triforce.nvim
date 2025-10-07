BEGIN;

ALTER TABLE unit_kerja
    ALTER COLUMN is_satker DROP DEFAULT;

ALTER TABLE unit_kerja
    ALTER COLUMN is_satker TYPE boolean USING (is_satker <> 0);

ALTER TABLE unit_kerja
    ALTER COLUMN is_satker SET NOT NULL,
    ALTER COLUMN is_satker SET DEFAULT false;

COMMIT;