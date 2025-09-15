-- name: ListRefJenisPenghargaan :many
SELECT 
  id, 
  nama 
FROM ref_jenis_penghargaan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisPenghargaan :one
SELECT COUNT(1) 
FROM ref_jenis_penghargaan
WHERE deleted_at IS NULL;
