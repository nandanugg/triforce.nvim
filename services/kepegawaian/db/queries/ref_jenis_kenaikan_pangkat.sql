-- name: ListJenisKenaikanPangkat :many
SELECT id, nama FROM ref_jenis_kenaikan_pangkat
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountJenisKenaikanPangkat :one
SELECT COUNT(1) FROM ref_jenis_kenaikan_pangkat
WHERE deleted_at IS NULL;

-- name: GetJenisKenaikanPangkat :one
SELECT 
    id,
    nama
FROM ref_jenis_kenaikan_pangkat
WHERE id = @id
  AND deleted_at IS NULL;

-- name: CreateJenisKenaikanPangkat :one
INSERT INTO ref_jenis_kenaikan_pangkat (
    nama
) VALUES (
    @nama
)
RETURNING id, nama;

-- name: UpdateJenisKenaikanPangkat :one
UPDATE ref_jenis_kenaikan_pangkat
SET 
    nama = @nama,
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING id, nama;

-- name: DeleteJenisKenaikanPangkat :execrows
UPDATE ref_jenis_kenaikan_pangkat
SET 
    deleted_at = now(),
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL;
