-- name: ListRefJenisJabatan :many
SELECT id, nama FROM ref_jenis_jabatan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisJabatan :one
SELECT COUNT(1) FROM ref_jenis_jabatan
WHERE deleted_at IS NULL;