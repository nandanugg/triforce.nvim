-- name: ListRefGolongan :many
SELECT id, nama, nama_pangkat, nama_2, gol_pppk, gol, created_at, updated_at FROM ref_golongan
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefGolongan :one
SELECT COUNT(1) FROM ref_golongan
WHERE deleted_at IS NULL;

-- name: GetRefGolongan :one
SELECT id, nama, nama_pangkat, nama_2, gol, gol_pppk, created_at, updated_at
FROM ref_golongan
WHERE id = @id::integer AND deleted_at IS NULL;

-- name: CreateRefGolongan :one
INSERT INTO ref_golongan (nama, nama_pangkat, nama_2, gol, gol_pppk)
VALUES (@nama::text, @nama_pangkat::text, @nama_2::text, @gol::smallint, @gol_pppk::text)
RETURNING id, nama, nama_pangkat, nama_2, gol, gol_pppk, created_at, updated_at;

-- name: UpdateRefGolongan :one
UPDATE ref_golongan
SET nama = @nama::text, nama_pangkat = @nama_pangkat::text, nama_2 = @nama_2::text, gol = @gol::smallint, gol_pppk = @gol_pppk::text, updated_at = NOW()
WHERE id = @id::integer AND deleted_at IS NULL
RETURNING id, nama, nama_pangkat, nama_2, gol, gol_pppk, created_at, updated_at;

-- name: DeleteRefGolongan :execrows
UPDATE ref_golongan
SET deleted_at = NOW()
WHERE id = @id::integer AND deleted_at IS NULL;
