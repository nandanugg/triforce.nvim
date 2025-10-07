BEGIN;

ALTER TABLE unit_kerja
    ALTER COLUMN is_satker DROP DEFAULT;

ALTER TABLE unit_kerja
    ALTER COLUMN is_satker TYPE int2 USING CASE WHEN is_satker THEN 1 ELSE 0 END;

ALTER TABLE unit_kerja
    ALTER COLUMN is_satker SET NOT NULL,
    ALTER COLUMN is_satker SET DEFAULT 0;

END;