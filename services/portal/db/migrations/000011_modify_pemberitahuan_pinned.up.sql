BEGIN;

ALTER TABLE pemberitahuan
	DROP COLUMN pinned,
        ADD COLUMN pinned_at TIMESTAMPTZ,
	ADD COLUMN aktif_range tstzrange
	GENERATED ALWAYS AS (tstzrange(diterbitkan_pada, ditarik_pada)) STORED;

CREATE INDEX idx_pemberitahuan_aktif_range
ON pemberitahuan USING gist (aktif_range)
WHERE deleted_at IS NULL;

COMMIT;
