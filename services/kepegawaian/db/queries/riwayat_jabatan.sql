-- name: ListRiwayatJabatan :many
SELECT 
    riwayat_jabatan.id,
    ref_jenis_jabatan.nama as jenis_jabatan,
    ref_jabatan.nama_jabatan,
    tmt_jabatan,
    no_sk,
    tanggal_sk,
    unit_kerja.nama_unor as satuan_kerja,
    status_plt,
    ref_kelas_jabatan.kelas_jabatan,
    periode_jabatan_start_date,
    periode_jabatan_end_date,
    unit_organisasi.nama_unor as unit_organisasi
FROM riwayat_jabatan 
JOIN unit_kerja on riwayat_jabatan.satuan_kerja_id = unit_kerja.id
JOIN ref_kelas_jabatan on riwayat_jabatan.kelas_jabatan_id = ref_kelas_jabatan.id
JOIN unit_kerja unit_organisasi on riwayat_jabatan.unor_id = unit_organisasi.id
JOIN ref_jenis_jabatan on riwayat_jabatan.jenis_jabatan_id = ref_jenis_jabatan.id
JOIN ref_jabatan on riwayat_jabatan.jabatan_id = ref_jabatan.id
WHERE riwayat_jabatan.pns_nip = @pns_nip::varchar and riwayat_jabatan.deleted_at IS NULL
ORDER BY tmt_jabatan DESC
LIMIT $1 OFFSET $2;

-- name: CountRiwayatJabatan :one
SELECT count(1)
FROM riwayat_jabatan 
JOIN unit_kerja on riwayat_jabatan.satuan_kerja_id = unit_kerja.id
JOIN ref_kelas_jabatan on riwayat_jabatan.kelas_jabatan_id = ref_kelas_jabatan.id
JOIN unit_kerja unit_organisasi on riwayat_jabatan.unor_id = unit_organisasi.id
JOIN ref_jenis_jabatan on riwayat_jabatan.jenis_jabatan_id = ref_jenis_jabatan.id
JOIN ref_jabatan on riwayat_jabatan.jabatan_id = ref_jabatan.id
WHERE riwayat_jabatan.pns_nip = @pns_nip::varchar and riwayat_jabatan.deleted_at IS NULL;