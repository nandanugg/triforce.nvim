-- name: ListSKByNIP :many
SELECT
    fds.file_id,
    fds.kategori as kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    fds.status_sk
FROM file_digital_signature fds
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::varchar
    AND (sqlc.narg('no_sk')::varchar IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::varchar || '%')
    AND (sqlc.narg('status_sk')::integer IS NULL OR fds.status_sk = sqlc.narg('status_sk')::integer)
    AND (sqlc.narg('kategori_sk')::varchar is null OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::varchar || '%')
ORDER BY fds.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSKByNIP :one
SELECT COUNT(1) as total
FROM file_digital_signature fds
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND (sqlc.narg('no_sk')::varchar IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::varchar || '%')
    AND (sqlc.narg('status_sk')::integer IS NULL OR fds.status_sk = sqlc.narg('status_sk')::integer)
    AND (sqlc.narg('kategori_sk')::varchar is null OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::varchar || '%');

-- name: GetSKByNIPAndID :one
SELECT
    fds.kategori as kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    fds.status_sk,
    p.nama as nama_pemilik_sk,
    pemroses.nama as nama_penandatangan
FROM file_digital_signature fds
JOIN pegawai p on p.nip_baru = fds.nip_sk  and p.deleted_at is null
LEFT JOIN pegawai pemroses on pemroses.nip_baru = fds.nip_pemroses and pemroses.deleted_at is null
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSKByNIPAndID :one
SELECT 
    file_base64
FROM 
    file_digital_signature fds
WHERE 
    fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSKSignedByNIPAndID :one
SELECT 
    file_base64_sign
FROM 
    file_digital_signature fds
WHERE 
    fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;