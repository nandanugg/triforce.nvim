BEGIN;

DROP INDEX IF EXISTS idx_pemberitahuan_aktif_range;

ALTER TABLE pemberitahuan
	DROP COLUMN pinned_at,
	ADD COLUMN pinned boolean NOT NULL DEFAULT false,
	DROP COLUMN IF EXISTS aktif_range;

COMMIT;
