-- name: ListRefTingkatPendidikan :many
SELECT
    tp.id,
    tp.nama,
    tp.golongan_id,
    tp.golongan_awal_id,
    tp.abbreviation,
    tp.tingkat
FROM ref_tingkat_pendidikan tp
WHERE tp.deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefTingkatPendidikan :one
SELECT COUNT(1)
FROM ref_tingkat_pendidikan
WHERE deleted_at IS NULL;

-- name: GetRefTingkatPendidikan :one
SELECT
    tp.id,
    tp.nama,
    tp.golongan_id,
    tp.golongan_awal_id,
    tp.abbreviation,
    tp.tingkat
FROM ref_tingkat_pendidikan tp
WHERE tp.deleted_at IS NULL AND tp.id = @id::integer;

-- name: CreateRefTingkatPendidikan :one
INSERT INTO ref_tingkat_pendidikan (nama, abbreviation, golongan_id, golongan_awal_id, tingkat)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat;

-- name: UpdateRefTingkatPendidikan :one
UPDATE ref_tingkat_pendidikan
SET nama = $1, abbreviation = $2, golongan_id = $3, golongan_awal_id = $4, tingkat = $5
WHERE id = $6 AND deleted_at IS NULL
RETURNING id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat;

-- name: DeleteRefTingkatPendidikan :execrows
UPDATE ref_tingkat_pendidikan
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
