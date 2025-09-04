-- name: GetRefJenisJabatan :many
SELECT id, nama FROM ref_jenis_jabatan
LIMIT $1 OFFSET $2;