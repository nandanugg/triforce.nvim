-- name: GetUnitKerjaByNamaOrInduk :many
SELECT id, nama_unor 
from unit_kerja
WHERE 
    (CASE WHEN @nama::varchar = '' THEN true ELSE nama_unor ilike @nama::varchar || '%' END)
    AND (CASE WHEN @unor_induk::varchar = '' THEN true ELSE unor_induk = @unor_induk::varchar END)
    AND deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountUnitKerja :one
SELECT COUNT(1) FROM unit_kerja
WHERE 
    (CASE WHEN @nama::varchar = '' THEN true ELSE nama_unor ilike @nama::varchar || '%' END)
    AND (CASE WHEN @unor_induk::varchar = '' THEN true ELSE unor_induk = @unor_induk::varchar END)
    AND deleted_at IS NULL;