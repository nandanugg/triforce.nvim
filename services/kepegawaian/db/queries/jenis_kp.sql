-- name: GetJenisKP :many
SELECT id, nama FROM ref_jenis_kp
LIMIT $1 OFFSET $2;