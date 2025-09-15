-- name: ListPelatihanTeknis :many
SELECT
    rk.id,
    rk.tipe_kursus,
    rk.jenis_kursus,
    rk.nama_kursus,
    rk.tanggal_kursus,
    EXTRACT(YEAR FROM rk.tanggal_kursus) as tahun,
    COALESCE(rk.lama_kursus, 0) as durasi,
    rk.institusi_penyelenggara,
    rk.no_sertifikat
FROM riwayat_kursus rk
JOIN pegawai p ON rk.pns_id = p.pns_id
WHERE p.nip_baru = $1
  AND p.deleted_at IS NULL
  AND rk.deleted_at IS NULL
ORDER BY rk.tanggal_kursus DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountPelatihanTeknis :one
SELECT COUNT(1)
FROM riwayat_kursus rk
JOIN pegawai p ON rk.pns_id = p.pns_id
WHERE p.nip_baru = $1
  AND p.deleted_at IS NULL
  AND rk.deleted_at IS NULL;
