-- name: ListRefAgama :many
SELECT id, nama, created_at, updated_at
FROM ref_agama
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CountRefAgama :one
SELECT COUNT(1)
FROM ref_agama
WHERE deleted_at IS NULL;

-- name: GetRefAgama :one
SELECT id, nama, created_at, updated_at
FROM ref_agama
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateRefAgama :one
INSERT INTO ref_agama (nama)
VALUES ($1)
RETURNING id, nama, created_at, updated_at;

-- name: UpdateRefAgama :one
UPDATE ref_agama
SET nama = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, nama, created_at, updated_at;

-- name: DeleteRefAgama :execrows
UPDATE ref_agama
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;
