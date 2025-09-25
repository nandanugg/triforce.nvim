-- name: ListSuratKeputusanByNIP :many
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

-- name: CountSuratKeputusanByNIP :one
SELECT COUNT(1) as total
FROM file_digital_signature fds
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND (sqlc.narg('no_sk')::varchar IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::varchar || '%')
    AND (sqlc.narg('status_sk')::integer IS NULL OR fds.status_sk = sqlc.narg('status_sk')::integer)
    AND (sqlc.narg('kategori_sk')::varchar is null OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::varchar || '%');

-- name: GetSuratKeputusanByNIPAndID :one
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

-- name: GetBerkasSuratKeputusanByNIPAndID :one
SELECT 
    file_base64
FROM 
    file_digital_signature fds
WHERE 
    fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanSignedByNIPAndID :one
SELECT 
    file_base64_sign
FROM 
    file_digital_signature fds
WHERE 
    fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: ListLogSuratKeputusanByID :many
SELECT 
    tindakan as log, 
    pemroses.nama as actor,
    fdsr.created_at as waktu_tindakan
FROM 
    file_digital_signature_riwayat fdsr
LEFT JOIN pegawai pemroses on pemroses.nip_baru = fdsr.nip_pemroses and pemroses.deleted_at is null
WHERE fdsr.file_id = @id::varchar and fdsr.deleted_at IS NULL;

-- name: ListSuratKeputusan :many
SELECT
    fds.file_id,
    p.nama as nama_pemilik_sk,
    fds.kategori AS kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    p.unor_id,
    fds.status_sk
FROM file_digital_signature fds
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fds.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('tanggal_sk_mulai')::DATE IS NULL OR fds.tanggal_sk >= sqlc.narg('tanggal_sk_mulai')::DATE)
    AND (sqlc.narg('tanggal_sk_akhir')::DATE IS NULL OR fds.tanggal_sk <= sqlc.narg('tanggal_sk_akhir')::DATE)
    AND (sqlc.narg('status_sk')::INTEGER IS NULL OR fds.status_sk = sqlc.narg('status_sk')::integer)
ORDER BY fds.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSuratKeputusan :one
SELECT COUNT(*) as total
FROM file_digital_signature fds
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fds.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('tanggal_sk_mulai')::DATE IS NULL OR fds.tanggal_sk >= sqlc.narg('tanggal_sk_mulai')::DATE)
    AND (sqlc.narg('tanggal_sk_akhir')::DATE IS NULL OR fds.tanggal_sk <= sqlc.narg('tanggal_sk_akhir')::DATE)
    AND (sqlc.narg('status_sk')::INTEGER IS NULL OR fds.status_sk = sqlc.narg('status_sk')::integer);

-- name: GetSuratKeputusanByID :one
SELECT
    fds.kategori as kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    fds.status_sk,
    fds.nip_sk,
    p.nama as nama_pemilik_sk,
    pemroses.nama as nama_penandatangan
FROM file_digital_signature fds
JOIN pegawai p on p.nip_baru = fds.nip_sk and p.deleted_at is null
LEFT JOIN pegawai pemroses on pemroses.nip_baru = fds.nip_pemroses and pemroses.deleted_at is null
WHERE fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanByID :one
SELECT 
    file_base64
FROM 
    file_digital_signature fds
WHERE 
    fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanSignedByID :one
SELECT 
    file_base64_sign
FROM 
    file_digital_signature fds
WHERE 
    fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar;
