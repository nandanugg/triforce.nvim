-- name: ListRiwayatPenghargaan :many
SELECT
    id,
    jenis_penghargaan,
    nama_penghargaan,
    deskripsi_penghargaan,
    tanggal_penghargaan
FROM riwayat_penghargaan_umum
WHERE nip = @nip::varchar and riwayat_penghargaan_umum.deleted_at is null
ORDER BY tanggal_penghargaan DESC NULLS LAST
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

-- name: CreateRiwayatPenghargaan :one
insert into 
    riwayat_penghargaan_umum (nip, nama_penghargaan, jenis_penghargaan, deskripsi_penghargaan, tanggal_penghargaan) 
values 
    ($1, $2, $3, $4, $5)
returning id;

-- name: UpdateRiwayatPenghargaan :execrows
update riwayat_penghargaan_umum
set
    nip = $2,
    nama_penghargaan = $3,
    jenis_penghargaan = $4,
    deskripsi_penghargaan = $5,
    tanggal_penghargaan = $6,
    updated_at = now()
where id = $1 AND deleted_at IS NULL;

-- name: UpdateRiwayatPenghargaanBerkas :execrows
update riwayat_penghargaan_umum
set
    file_base64 = @file_base64,
    updated_at = now()
where id = @id AND nip = @nip::varchar AND riwayat_penghargaan_umum.deleted_at IS NULL;

-- name: DeleteRiwayatPenghargaan :execrows
update riwayat_penghargaan_umum
set
    deleted_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;
