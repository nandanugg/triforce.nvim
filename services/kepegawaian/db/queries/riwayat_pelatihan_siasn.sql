-- name: ListRiwayatPelatihanSIASN :many
SELECT
    rd.id,
    rjd.jenis_diklat,
    rd.nama_diklat,
    rd.no_sertifikat,
    rd.tanggal_mulai,
    rd.tanggal_selesai,
    rd.tahun_diklat,
    rd.durasi_jam,
    rd.institusi_penyelenggara
FROM riwayat_diklat rd
LEFT JOIN ref_jenis_diklat rjd ON rd.jenis_diklat_id = rjd.id AND rjd.deleted_at IS NULL
WHERE rd.nip_baru = $1 AND rd.deleted_at IS NULL
ORDER BY rd.tanggal_selesai DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanSIASN :one
SELECT count(*) FROM riwayat_diklat
WHERE nip_baru = $1 AND deleted_at IS NULL;

-- name: GetBerkasRiwayatPelatihanSIASN :one
select file_base64 from riwayat_diklat
where nip_baru = $1 and id = $2 and deleted_at is null;
