-- name: GetRefJabatan :many
select id, kode_jabatan, nama_jabatan from ref_jabatan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJabatan :one
SELECT COUNT(1) FROM ref_jabatan
WHERE deleted_at IS NULL;