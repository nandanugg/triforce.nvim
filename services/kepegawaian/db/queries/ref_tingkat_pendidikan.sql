-- name: ListRefTingkatPendidikan :many
SELECT
    tp.id,
    tp.nama
FROM ref_tingkat_pendidikan tp
WHERE tp.deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefTingkatPendidikan :one
SELECT COUNT(1)
FROM ref_tingkat_pendidikan
WHERE deleted_at IS NULL;
