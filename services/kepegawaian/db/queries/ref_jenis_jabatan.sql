-- name: CreateRefJenisJabatan :one
INSERT INTO ref_jenis_jabatan (nama)
VALUES ($1)
RETURNING id, nama;

-- name: UpdateRefJenisJabatan :one
UPDATE ref_jenis_jabatan
SET nama = $2
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, nama;

-- name: DeleteRefJenisJabatan :execrows
UPDATE ref_jenis_jabatan
SET deleted_at = now()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetRefJenisJabatan :one
SELECT id, nama
FROM ref_jenis_jabatan
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListRefJenisJabatan :many
SELECT id, nama FROM ref_jenis_jabatan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisJabatan :one
SELECT COUNT(1) FROM ref_jenis_jabatan
WHERE deleted_at IS NULL;