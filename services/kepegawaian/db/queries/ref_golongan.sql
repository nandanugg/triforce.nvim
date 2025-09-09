-- name: GetRefGolongan :many
SELECT id, nama, nama_pangkat FROM ref_golongan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefGolongan :one
SELECT COUNT(1) FROM ref_golongan
WHERE deleted_at IS NULL;