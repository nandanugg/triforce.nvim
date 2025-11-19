-- name: CreateRefJenisHukuman :one
INSERT INTO ref_jenis_hukuman (nama, tingkat_hukuman)
VALUES (@nama, @tingkat)
RETURNING id, nama, tingkat_hukuman as tingkat;

-- name: UpdateRefJenisHukuman :one
UPDATE ref_jenis_hukuman
SET 
  nama = @nama,
  tingkat_hukuman = @tingkat
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, nama, tingkat_hukuman as tingkat;

-- name: DeleteRefJenisHukuman :execrows
UPDATE ref_jenis_hukuman
SET deleted_at = now()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetRefJenisHukuman :one
SELECT id, nama, tingkat_hukuman as tingkat
FROM ref_jenis_hukuman
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListRefJenisHukuman :many
SELECT 
  id, 
  nama,
  tingkat_hukuman as tingkat
FROM ref_jenis_hukuman
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJenisHukuman :one
SELECT COUNT(1) 
FROM ref_jenis_hukuman
WHERE deleted_at IS NULL;

-- name: IsExistReferencesRiwayatHukumanDisiplinByID :one
SELECT EXISTS (
    SELECT 1
    FROM riwayat_hukuman_disiplin
    WHERE jenis_hukuman_id = $1::int
      AND deleted_at IS NULL
) AS exists;