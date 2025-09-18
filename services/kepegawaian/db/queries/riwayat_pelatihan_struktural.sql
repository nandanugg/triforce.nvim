-- name: ListRiwayatPelatihanStruktural :many
SELECT
	rjd.nama as jenis_diklat,
	rds.id,
	rds.nama_diklat,
	rds.nomor,
	rds.tanggal,
	rds.tahun,
	rds.institusi_penyelenggara,
	rds.lama
FROM riwayat_diklat_struktural rds
LEFT JOIN ref_jenis_diklat_struktural rjd ON rds.jenis_diklat_id = rjd.id AND rjd.deleted_at IS NULL
WHERE rds.pns_nip = $1 AND rds.deleted_at IS NULL
ORDER BY rds.tanggal DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanStruktural :one
SELECT count(*) FROM riwayat_diklat_struktural
WHERE pns_nip = $1 AND deleted_at IS NULL;
