-- name: ListRefJenisDiklatFungsional :many
SELECT 
  id, 
  nama 
FROM ref_jenis_diklat_fungsional
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisDiklatFungsional :one
SELECT COUNT(1) 
FROM ref_jenis_diklat_fungsional
WHERE deleted_at IS NULL;

