-- name: ListRefJabatan :many
SELECT kode_jabatan, id, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id, tunjangan_jabatan, created_at, updated_at
FROM ref_jabatan
WHERE (sqlc.narg('nama')::varchar IS NULL OR nama_jabatan ILIKE sqlc.narg('nama')::varchar || '%')
  AND deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJabatan :one
SELECT COUNT(1) FROM ref_jabatan
WHERE (sqlc.narg('nama')::varchar IS NULL OR nama_jabatan ILIKE sqlc.narg('nama')::varchar || '%')
  AND deleted_at IS NULL;

-- name: ListRefJabatanWithKeyword :many
SELECT j.kode_jabatan, j.id, j.nama_jabatan, j.nama_jabatan_full, j.jenis_jabatan, rj.nama as jenis_jabatan_nama, j.kelas, j.pensiun, j.kode_bkn, j.nama_jabatan_bkn, j.kategori_jabatan, j.bkn_id, j.tunjangan_jabatan, j.created_at, j.updated_at
FROM ref_jabatan j
LEFT JOIN ref_jenis_jabatan rj ON rj.id = j.jenis_jabatan and rj.deleted_at IS NULL
WHERE
  (sqlc.narg('keyword')::varchar IS NULL OR nama_jabatan ILIKE '%' || sqlc.narg('keyword')::varchar || '%' OR kategori_jabatan ILIKE '%' || sqlc.narg('keyword')::varchar || '%')
  AND j.deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountRefJabatanWithKeyword :one
SELECT COUNT(1) FROM ref_jabatan
WHERE
  (sqlc.narg('keyword')::varchar IS NULL OR nama_jabatan ILIKE '%' || sqlc.narg('keyword')::varchar || '%' OR kategori_jabatan ILIKE '%' || sqlc.narg('keyword')::varchar || '%')
  AND deleted_at IS NULL;

-- name: GetRefJabatan :one
SELECT kode_jabatan, id, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id, tunjangan_jabatan, created_at, updated_at
FROM ref_jabatan
WHERE id = @id::int AND deleted_at IS NULL;

-- name: GetRefJabatanByKode :one
select
  kode_jabatan as kode,
  nama_jabatan as nama,
  jenis_jabatan as jenis,
  kelas,
  kode_bkn
from ref_jabatan
where kode_jabatan = $1 and deleted_at is null;

-- name: CreateRefJabatan :one
INSERT INTO
  ref_jabatan (kode_jabatan, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id, tunjangan_jabatan)
VALUES
  (@kode_jabatan, @nama_jabatan, @nama_jabatan_full, @jenis_jabatan, @kelas, @pensiun, @kode_bkn, @nama_jabatan_bkn, @kategori_jabatan, @bkn_id, @tunjangan_jabatan)
RETURNING
  id, kode_jabatan, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id, tunjangan_jabatan, created_at, updated_at;

-- name: UpdateRefJabatan :one
UPDATE ref_jabatan
SET
  nama_jabatan = @nama_jabatan,
  nama_jabatan_full = @nama_jabatan_full,
  jenis_jabatan = @jenis_jabatan,
  kelas = @kelas,
  pensiun = @pensiun,
  kode_bkn = @kode_bkn,
  nama_jabatan_bkn = @nama_jabatan_bkn,
  kategori_jabatan = @kategori_jabatan,
  bkn_id = @bkn_id,
  updated_at = NOW(),
  kode_jabatan = @kode_jabatan,
  tunjangan_jabatan = @tunjangan_jabatan
WHERE
  id = @id::int AND deleted_at IS NULL
RETURNING
  kode_jabatan, id, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id, tunjangan_jabatan, created_at, updated_at;

-- name: DeleteRefJabatan :execrows
UPDATE ref_jabatan
SET deleted_at = NOW()
WHERE id = @id::int AND deleted_at IS NULL;

-- name: IsExistReferencesPegawaiByID :one
SELECT EXISTS (
    SELECT 1
    FROM pegawai
    JOIN ref_jabatan 
      ON ref_jabatan.kode_jabatan = pegawai.jabatan_instansi_id
     AND ref_jabatan.deleted_at IS NULL
    WHERE ref_jabatan.id = @id::int
      AND pegawai.deleted_at IS NULL
) AS exists;
