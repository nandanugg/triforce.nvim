BEGIN;

ALTER TABLE pemberitahuan
	ADD COLUMN pinned boolean NOT NULL DEFAULT false,
        ADD COLUMN diterbitkan_pada TIMESTAMPTZ NOT NULL,
        ADD COLUMN ditarik_pada TIMESTAMPTZ NOT NULL,
        ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        ADD COLUMN deleted_at TIMESTAMPTZ,
        ALTER COLUMN updated_at SET DEFAULT NOW(),
	DROP COLUMN status;

COMMIT;
