-- name: ListRiwayatJabatan :many
SELECT
    riwayat_jabatan.id,
    jenis_jabatan_id,
    ref_jenis_jabatan.nama as jenis_jabatan,
    riwayat_jabatan.jabatan_id as id_jabatan,
    ref_jabatan.nama_jabatan,
    tmt_jabatan,
    no_sk,
    tanggal_sk,
    satuan_kerja_id,
    ref_unit_kerja.nama_unor as satuan_kerja,
    status_plt,
    kelas_jabatan_id,
    ref_kelas_jabatan.kelas_jabatan,
    periode_jabatan_start_date,
    periode_jabatan_end_date,
    unor_id as unit_organisasi_id,
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

-- name: CreateRiwayatJabatan :one
insert into riwayat_jabatan
    (jenis_jabatan_id, jenis_jabatan, jabatan_id, nama_jabatan, jabatan_id_bkn, satuan_kerja_id, unor_id, unor_id_bkn, unor, tmt_jabatan, no_sk, tanggal_sk, status_plt, periode_jabatan_start_date, periode_jabatan_end_date, pns_id, pns_nip, pns_nama) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
returning id;

-- name: UpdateRiwayatJabatan :execrows
update riwayat_jabatan
set
    jenis_jabatan_id = $1,
    jenis_jabatan = $2,
    jabatan_id = $3,
    nama_jabatan = $4,
    jabatan_id_bkn = $5,
    satuan_kerja_id = $6,
    unor_id = $7,
    unor_id_bkn = $8,
    unor = $9,
    tmt_jabatan = $10,
    no_sk = $11,
    tanggal_sk = $12,
    status_plt = $13,
    periode_jabatan_start_date = $14,
    periode_jabatan_end_date = $15,
    updated_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatJabatan :execrows
update riwayat_jabatan
set deleted_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: UploadBerkasRiwayatJabatan :execrows
update riwayat_jabatan
set
    file_base64 = $1,
    updated_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;
