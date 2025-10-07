-- name: CreateRefJenisSatker :one
INSERT INTO ref_jenis_satker (nama)
VALUES ($1)
RETURNING id, nama;

-- name: UpdateRefJenisSatker :one
UPDATE ref_jenis_satker
SET nama = $2
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, nama;

-- name: DeleteRefJenisSatker :execrows
UPDATE ref_jenis_satker
SET deleted_at = now()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetRefJenisSatker :one
SELECT id, nama
FROM ref_jenis_satker
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListRefJenisSatker :many
SELECT id, nama
FROM ref_jenis_satker
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama')::varchar IS NULL OR nama ILIKE '%' || sqlc.narg('nama') || '%')
LIMIT $1 OFFSET $2;

-- name: CountRefJenisSatker :one
SELECT COUNT(*) AS total
FROM ref_jenis_satker
WHERE deleted_at IS NULL
  AND (sqlc.narg('nama')::varchar IS NULL OR nama ILIKE '%' || sqlc.narg('nama') || '%');