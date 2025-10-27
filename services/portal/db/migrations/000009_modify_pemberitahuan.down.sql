BEGIN;

ALTER TABLE pemberitahuan
    DROP COLUMN pinned,
    DROP COLUMN diterbitkan_pada,
    DROP COLUMN ditarik_pada,
    DROP COLUMN created_at,
    DROP COLUMN deleted_at,
    ALTER COLUMN updated_at DROP DEFAULT,
    ADD COLUMN status TEXT;

COMMIT;
