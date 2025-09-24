-- name: ListRiwayatPelatihanFungsional :many
SELECT
    rf.id,
    rf.jenis_diklat,
    rf.nama_kursus,
    rf.tanggal_kursus,
    rf.tahun,
    rf.institusi_penyelenggara,
    rf.no_sertifikat,
    rf.jumlah_jam
FROM riwayat_diklat_fungsional rf
WHERE rf.nip_baru = $1
AND rf.deleted_at IS NULL
ORDER BY rf."tahun" DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanFungsional :one
SELECT COUNT(1) AS total
FROM riwayat_diklat_fungsional rf
WHERE rf.nip_baru = $1
  AND rf.deleted_at IS NULL;

-- name: GetBerkasRiwayatPelatihanFungsional :one
select file_base64 from riwayat_diklat_fungsional
where nip_baru = $1 and id = $2 and deleted_at is null;
