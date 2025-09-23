-- name: ListRiwayatPenghargaan :many
SELECT
    id,
    jenis_penghargaan,
    nama_penghargaan,
    deskripsi_penghargaan,
    tanggal_penghargaan
FROM riwayat_penghargaan_umum
WHERE nip = @nip::varchar and riwayat_penghargaan_umum.deleted_at is null
ORDER BY tanggal_penghargaan DESC
LIMIT $1 OFFSET $2;

-- name: CountRiwayatPenghargaan :one
SELECT COUNT(1)
FROM riwayat_penghargaan_umum
WHERE nip = @nip::varchar and riwayat_penghargaan_umum.deleted_at is null;

-- name: GetBerkasRiwayatPenghargaan :one
SELECT file_base64
FROM riwayat_penghargaan_umum rpu
WHERE nip = $1
  AND rpu.id = $2
  AND rpu.deleted_at is null;
