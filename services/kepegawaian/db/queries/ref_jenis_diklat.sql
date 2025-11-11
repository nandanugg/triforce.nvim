-- name: ListRefJenisDiklat :many
SELECT
  id,
  jenis_diklat
FROM ref_jenis_diklat
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisDiklat :one
SELECT COUNT(1)
FROM ref_jenis_diklat
WHERE deleted_at IS NULL;

-- name: GetRefJenisDiklat :one
select id, jenis_diklat
from ref_jenis_diklat
where id = $1 and deleted_at is null;
