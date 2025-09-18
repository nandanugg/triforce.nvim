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
JOIN pegawai p ON rk.pns_id = p.pns_id AND p.deleted_at IS NULL
WHERE p.nip_baru = $1 AND rk.deleted_at IS NULL
ORDER BY rk.tanggal_kursus DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanTeknis :one
SELECT COUNT(1)
FROM riwayat_kursus rk
JOIN pegawai p ON rk.pns_id = p.pns_id AND p.deleted_at IS NULL
WHERE p.nip_baru = $1 AND rk.deleted_at IS NULL;
