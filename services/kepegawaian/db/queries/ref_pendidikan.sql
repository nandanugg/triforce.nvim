-- name: ListRefPendidikanWithTingkatPendidikan :many
SELECT
    p.id, p.nama, tp.nama as tingkat_pendidikan, tingkat_pendidikan_id
FROM
    ref_pendidikan p
LEFT JOIN
    ref_tingkat_pendidikan tp ON tp.id = p.tingkat_pendidikan_id AND tp.deleted_at IS NULL
WHERE
    p.deleted_at IS NULL
    AND (sqlc.narg('nama')::varchar IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama')::varchar || '%')
LIMIT $1 OFFSET $2;

-- name: CountRefPendidikan :one
SELECT
    COUNT(1)
FROM
    ref_pendidikan p
WHERE
    deleted_at IS NULL
    AND (sqlc.narg('nama')::varchar IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama')::varchar || '%');

-- name: GetRefPendidikan :one
SELECT
    p.id, p.nama, tp.nama as tingkat_pendidikan, tingkat_pendidikan_id
FROM
    ref_pendidikan p
LEFT JOIN
    ref_tingkat_pendidikan tp ON tp.id = p.tingkat_pendidikan_id AND tp.deleted_at IS NULL
WHERE
    p.id = @id::text AND p.deleted_at IS NULL;

-- name: CreateRefPendidikan :one
INSERT INTO
    ref_pendidikan (id, nama, tingkat_pendidikan_id)
VALUES
    (@id::text, @nama::text, @tingkat_pendidikan_id)
RETURNING id;

-- name: UpdateRefPendidikan :one
UPDATE ref_pendidikan
SET nama = @nama::text, tingkat_pendidikan_id = @tingkat_pendidikan_id
WHERE id = @id::text AND deleted_at IS NULL
RETURNING id;

-- name: DeleteRefPendidikan :execrows
UPDATE ref_pendidikan
SET deleted_at = NOW()
WHERE id = @id::text AND deleted_at IS NULL;
