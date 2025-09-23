-- name: ListRefJabatan :many
select kode_jabatan, nama_jabatan from ref_jabatan
WHERE (sqlc.narg('nama')::varchar IS NULL OR nama_jabatan ILIKE sqlc.narg('nama')::varchar || '%')
  AND deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJabatan :one
SELECT COUNT(1) FROM ref_jabatan
WHERE (sqlc.narg('nama')::varchar IS NULL OR nama_jabatan ILIKE sqlc.narg('nama')::varchar || '%')
  AND deleted_at IS NULL;
