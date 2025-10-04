-- name: ListRiwayatPelatihanStruktural :many
SELECT
    rs.id,
    rs.nama_diklat,
    rs.tanggal,
    rs.tahun,
    rs.nomor,
    rs.lama
FROM riwayat_diklat_struktural rs
WHERE rs.pns_nip = $1
AND rs.deleted_at IS NULL
ORDER BY rs."tahun" DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanStruktural :one
SELECT COUNT(1) AS total
FROM riwayat_diklat_struktural rs
WHERE rs.pns_nip = $1
  AND rs.deleted_at IS NULL;

-- name: GetBerkasRiwayatPelatihanStruktural :one
select file_base64 from riwayat_diklat_struktural
where pns_nip = $1 and id = $2 and deleted_at is null;
