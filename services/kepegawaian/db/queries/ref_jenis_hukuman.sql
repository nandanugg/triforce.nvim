-- name: ListRefJenisHukuman :many
SELECT 
  id, 
  nama 
FROM ref_jenis_hukuman
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisHukuman :one
SELECT COUNT(1) 
FROM ref_jenis_hukuman
WHERE deleted_at IS NULL;
