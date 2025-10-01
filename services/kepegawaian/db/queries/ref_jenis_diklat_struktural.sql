-- name: ListRefJenisDiklatStruktural :many
SELECT 
  id, 
  nama 
FROM ref_jenis_diklat_struktural
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisDiklatStruktural :one
SELECT COUNT(1) 
FROM ref_jenis_diklat_struktural
WHERE deleted_at IS NULL;

