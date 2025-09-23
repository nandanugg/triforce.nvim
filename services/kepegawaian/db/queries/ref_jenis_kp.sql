-- name: ListJenisKP :many
SELECT id, nama FROM ref_jenis_kenaikan_pangkat
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountJenisKP :one
SELECT COUNT(1) FROM ref_jenis_kenaikan_pangkat
WHERE deleted_at IS NULL;
