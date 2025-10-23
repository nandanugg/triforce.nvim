-- name: ListRiwayatJabatan :many
SELECT
    riwayat_jabatan.id,
    ref_jenis_jabatan.nama as jenis_jabatan,
    riwayat_jabatan.jabatan_id as id_jabatan,
    ref_jabatan.nama_jabatan,
    tmt_jabatan,
    no_sk,
    tanggal_sk,
    ref_unit_kerja.nama_unor as satuan_kerja,
    status_plt,
    ref_kelas_jabatan.kelas_jabatan,
    periode_jabatan_start_date,
    periode_jabatan_end_date,
    unit_organisasi.nama_unor as unit_organisasi
FROM riwayat_jabatan
LEFT JOIN ref_unit_kerja on riwayat_jabatan.satuan_kerja_id = ref_unit_kerja.id AND ref_unit_kerja.deleted_at IS NULL
LEFT JOIN ref_kelas_jabatan on riwayat_jabatan.kelas_jabatan_id = ref_kelas_jabatan.id AND ref_kelas_jabatan.deleted_at IS NULL
LEFT JOIN ref_unit_kerja unit_organisasi on riwayat_jabatan.unor_id = unit_organisasi.id AND unit_organisasi.deleted_at IS NULL
LEFT JOIN ref_jenis_jabatan on riwayat_jabatan.jenis_jabatan_id = ref_jenis_jabatan.id AND ref_jenis_jabatan.deleted_at IS NULL
LEFT JOIN ref_jabatan on riwayat_jabatan.jabatan_id = ref_jabatan.kode_jabatan AND ref_jabatan.deleted_at IS NULL
WHERE riwayat_jabatan.pns_nip = @pns_nip::varchar and riwayat_jabatan.deleted_at IS NULL
ORDER BY tmt_jabatan DESC
LIMIT $1 OFFSET $2;

-- name: CountRiwayatJabatan :one
SELECT count(1)
FROM riwayat_jabatan
WHERE riwayat_jabatan.pns_nip = @pns_nip::varchar and riwayat_jabatan.deleted_at IS NULL;

-- name: GetBerkasRiwayatJabatan :one
select file_base64 from riwayat_jabatan
where pns_nip = $1 and id = $2 and deleted_at is null;
