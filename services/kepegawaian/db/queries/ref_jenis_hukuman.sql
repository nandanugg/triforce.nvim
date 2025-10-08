-- name: CreateRefJenisHukuman :one
INSERT INTO ref_jenis_hukuman (nama)
VALUES ($1)
RETURNING id, nama;

-- name: UpdateRefJenisHukuman :one
UPDATE ref_jenis_hukuman
SET nama = $2
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, nama;

-- name: DeleteRefJenisHukuman :execrows
UPDATE ref_jenis_hukuman
SET deleted_at = now()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetRefJenisHukuman :one
SELECT id, nama
FROM ref_jenis_hukuman
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListRefJenisHukuman :many
SELECT 
  id, 
  nama 
FROM ref_jenis_hukuman
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisHukuman :one
SELECT COUNT(1) 
FROM ref_jenis_hukuman
WHERE deleted_at IS NULL;
