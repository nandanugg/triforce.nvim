-- name: CreateRefJenisPenghargaan :one
INSERT INTO ref_jenis_penghargaan (nama)
VALUES ($1)
RETURNING id, nama;

-- name: UpdateRefJenisPenghargaan :one
UPDATE ref_jenis_penghargaan
SET nama = $2
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, nama;

-- name: DeleteRefJenisPenghargaan :execrows
UPDATE ref_jenis_penghargaan
SET deleted_at = now()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetRefJenisPenghargaan :one
SELECT id, nama
FROM ref_jenis_penghargaan
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListRefJenisPenghargaan :many
SELECT 
  id, 
  nama 
FROM ref_jenis_penghargaan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisPenghargaan :one
SELECT COUNT(1) 
FROM ref_jenis_penghargaan
WHERE deleted_at IS NULL;
