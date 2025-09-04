-- name: GetRefGolongan :many
SELECT id, nama, nama_pangkat FROM ref_golongan
LIMIT $1 OFFSET $2;