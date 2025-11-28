-- name: ListRefLokasi :many
SELECT 
    p.id, p.nama
FROM 
    ref_lokasi p
WHERE 
    p.deleted_at IS NULL
    AND (@nama::varchar = '' OR p.nama ILIKE '%' || @nama::varchar || '%')
LIMIT $1 OFFSET $2;

-- name: CountRefLokasi :one
SELECT COUNT(1)
FROM ref_lokasi p
WHERE 
    p.deleted_at IS NULL
    AND (@nama::varchar = '' OR p.nama ILIKE '%' || @nama::varchar || '%');
