-- name: ListRiwayatPelatihanTeknis :many
SELECT
    rk.id,
    rk.tipe_kursus,
    rk.jenis_kursus,
    rk.nama_kursus,
    rk.tanggal_kursus,
    rk.lama_kursus as durasi,
    rk.institusi_penyelenggara,
    rk.no_sertifikat
FROM riwayat_kursus rk
WHERE rk.pns_nip = $1 AND rk.deleted_at IS NULL
ORDER BY rk.tanggal_kursus DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanTeknis :one
SELECT COUNT(1)
FROM riwayat_kursus rk
WHERE rk.pns_nip = $1 AND rk.deleted_at IS NULL;

-- name: GetBerkasRiwayatPelatihanTeknis :one
select file_base64 from riwayat_kursus
where pns_nip = $1 and id = $2 and deleted_at is null;
