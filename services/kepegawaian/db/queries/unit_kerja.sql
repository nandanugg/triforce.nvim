-- name: ListUnitKerjaByNamaOrInduk :many
SELECT id, nama_unor 
from unit_kerja
WHERE 
    (CASE WHEN @search::varchar = '' THEN true ELSE nama_unor ilike @search::varchar || '%' END)
    AND (CASE WHEN @unor_induk::varchar = '' THEN true ELSE unor_induk = @unor_induk::varchar END)
LIMIT $1 OFFSET $2;