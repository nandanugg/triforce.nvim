-- name: ListRefTingkatPendidikan :many
SELECT
    tp.id,
    tp.nama
FROM ref_tingkat_pendidikan tp
WHERE tp.deleted_at IS NULL
;
