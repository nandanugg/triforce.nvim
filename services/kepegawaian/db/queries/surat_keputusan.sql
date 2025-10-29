-- name: ListSuratKeputusanByNIP :many
SELECT
    fds.file_id,
    fds.kategori as kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    fds.status_sk
FROM surat_keputusan fds
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::varchar
    AND (sqlc.narg('no_sk')::varchar IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::varchar || '%')
    AND (sqlc.narg('list_status_sk')::integer[] IS NULL OR fds.status_sk = ANY(sqlc.narg('list_status_sk')::integer[]))
    AND (sqlc.narg('kategori_sk')::varchar is null OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::varchar || '%')
    AND ds_ok = true
ORDER BY fds.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSuratKeputusanByNIP :one
SELECT COUNT(1) as total
FROM surat_keputusan fds
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND (sqlc.narg('no_sk')::varchar IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::varchar || '%')
    AND (sqlc.narg('list_status_sk')::integer[] IS NULL OR fds.status_sk = ANY(sqlc.narg('list_status_sk')::integer[]))
    AND (sqlc.narg('kategori_sk')::varchar is null OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::varchar || '%')
    AND ds_ok = true;

-- name: GetSuratKeputusanByNIPAndID :one
SELECT
    fds.kategori as kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    fds.status_sk,
    p.nama as nama_pemilik_sk,
    pemroses.nama as nama_penandatangan
FROM surat_keputusan fds
JOIN pegawai p on p.nip_baru = fds.nip_sk  and p.deleted_at is null
LEFT JOIN pegawai pemroses on pemroses.nip_baru = fds.nip_pemroses and pemroses.deleted_at is null
WHERE fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanByNIPAndID :one
SELECT
    file_base64
FROM
    surat_keputusan fds
WHERE
    fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanSignedByNIPAndID :one
SELECT
    file_base64_sign
FROM
    surat_keputusan fds
WHERE
    fds.deleted_at IS NULL
    AND fds.nip_sk = @nip::VARCHAR
    AND fds.file_id = @id::varchar;

-- name: ListLogSuratKeputusanByID :many
SELECT
    tindakan as log,
    pemroses.nama as actor,
    fdsr.created_at as waktu_tindakan
FROM
    riwayat_surat_keputusan fdsr
LEFT JOIN pegawai pemroses on pemroses.nip_baru = fdsr.nip_pemroses and pemroses.deleted_at is null
WHERE fdsr.file_id = @id::varchar and fdsr.deleted_at IS NULL;

-- name: ListSuratKeputusan :many
SELECT
    fds.file_id,
    p.nama as nama_pemilik_sk,
    p.nip_baru as nip_pemilik_sk,
    fds.kategori AS kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    p.unor_id,
    fds.status_sk
FROM surat_keputusan fds
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fds.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('tanggal_sk_mulai')::DATE IS NULL OR fds.tanggal_sk >= sqlc.narg('tanggal_sk_mulai')::DATE)
    AND (sqlc.narg('tanggal_sk_akhir')::DATE IS NULL OR fds.tanggal_sk <= sqlc.narg('tanggal_sk_akhir')::DATE)
    AND (sqlc.narg('list_status_sk')::integer[] IS NULL OR fds.status_sk = ANY(sqlc.narg('list_status_sk')::integer[]))
    AND ds_ok = true
ORDER BY fds.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSuratKeputusan :one
SELECT COUNT(*) as total
FROM surat_keputusan fds
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fds.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('tanggal_sk_mulai')::DATE IS NULL OR fds.tanggal_sk >= sqlc.narg('tanggal_sk_mulai')::DATE)
    AND (sqlc.narg('tanggal_sk_akhir')::DATE IS NULL OR fds.tanggal_sk <= sqlc.narg('tanggal_sk_akhir')::DATE)
    AND (sqlc.narg('list_status_sk')::integer[] IS NULL OR fds.status_sk = ANY(sqlc.narg('list_status_sk')::integer[]))
    AND ds_ok = true;

-- name: GetSuratKeputusanByID :one
SELECT
    fds.kategori as kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    fds.status_sk,
    p.nama as nama_pemilik_sk,
    p.nip_baru as nip_pemilik_sk,
    pemroses.nama as nama_penandatangan,
    rj.nama_jabatan as jabatan_penandatangan,
    status_ttd,
    status_koreksi,
    catatan,
    ttd_pegawai_id
FROM surat_keputusan fds
JOIN pegawai p on p.nip_baru = fds.nip_sk and p.deleted_at is null
LEFT JOIN pegawai pemroses on pemroses.nip_baru = fds.nip_pemroses and pemroses.deleted_at is null
LEFT JOIN ref_jabatan rj on pemroses.jabatan_instansi_id = rj.kode_jabatan and rj.deleted_at is null
WHERE fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanByID :one
SELECT
    file_base64
FROM
    surat_keputusan fds
WHERE
    fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar;

-- name: GetBerkasSuratKeputusanSignedByID :one
SELECT
    file_base64_sign
FROM
    surat_keputusan fds
WHERE
    fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar;

-- name: ListKoreksiSuratKeputusanByPNSID :many
SELECT
    fds.file_id,
    p.nama as nama_pemilik_sk,
    p.nip_baru as nip_pemilik_sk,
    fds.kategori AS kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    p.unor_id
FROM koreksi_surat_keputusan fdc
JOIN surat_keputusan fds ON fds.file_id = fdc.file_id AND fds.deleted_at IS NULL
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fdc.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip_pemilik')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip_pemilik')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('no_sk')::VARCHAR IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::VARCHAR || '%')
    AND (
        sqlc.narg('status_koreksi')::integer[] IS NULL
        OR fdc.status_koreksi = ANY(sqlc.narg('status_koreksi')::integer[])
    )
    AND fdc.pegawai_korektor_id = @pns_id::varchar
    AND ds_ok = true
ORDER BY fds.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountKoreksiSuratKeputusanByPNSID :one
select
    count(1) as total
FROM koreksi_surat_keputusan fdc
JOIN surat_keputusan fds ON fds.file_id = fdc.file_id AND fds.deleted_at IS NULL
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fdc.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip_pemilik')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip_pemilik')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('no_sk')::VARCHAR IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::VARCHAR || '%')
    AND (
        sqlc.narg('status_koreksi')::integer[] IS NULL
        OR fdc.status_koreksi = ANY(sqlc.narg('status_koreksi')::integer[])
    )
    AND fdc.pegawai_korektor_id = @pns_id::varchar
    AND ds_ok = true;

-- name: ListAntreanKoreksiSuratKeputusanByNIP :many
SELECT 
  fds.kategori,
  COUNT(1) AS jumlah
FROM 
  surat_keputusan fds
JOIN 
  koreksi_surat_keputusan fdc ON fds.file_id = fdc.file_id AND fdc.deleted_at IS NULL
JOIN 
  pegawai korektor ON fdc.pegawai_korektor_id = korektor.pns_id AND korektor.deleted_at IS NULL
WHERE 
  fds.status_ttd = 0
  AND fds.ds_ok = true
  AND fds.kategori != '< Semua >'
  AND (fdc.status_koreksi = 0 OR fdc.status_koreksi IS NULL)
  AND korektor.nip_baru = @nip_korektor::varchar
  AND fds.ds_ok = true
GROUP BY 
  fds.kategori
LIMIT $1 OFFSET $2;

-- name: CountAntreanKoreksiSuratKeputusanByNIP :one
SELECT 
    COUNT(DISTINCT fds.kategori) AS total
FROM 
    surat_keputusan fds
JOIN 
    koreksi_surat_keputusan fdc ON fds.file_id = fdc.file_id AND fdc.deleted_at IS NULL
JOIN 
    pegawai korektor ON fdc.pegawai_korektor_id = korektor.pns_id AND korektor.deleted_at IS NULL
WHERE 
    fds.status_ttd = 0
    AND fds.ds_ok = true
    AND fds.kategori != '< Semua >'
    AND (fdc.status_koreksi = 0 OR fdc.status_koreksi IS NULL)
    AND korektor.nip_baru = @nip_korektor::varchar
    AND fds.ds_ok = true;

-- name: ListKorektorSuratKeputusanByID :many
SELECT 
    fds.file_id,
    fdc.korektor_ke,
    korektor.nama as nama_korektor,
    korektor.nip_baru as nip_korektor,
    korektor.gelar_depan as gelar_depan_korektor,
    korektor.gelar_belakang as gelar_belakang_korektor,
    fdc.status_koreksi,
    fdc.catatan_koreksi,
    fdc.pegawai_korektor_id
FROM 
    surat_keputusan fds
JOIN 
    koreksi_surat_keputusan fdc ON fds.file_id = fdc.file_id AND fdc.deleted_at IS NULL
JOIN 
    pegawai korektor on fdc.pegawai_korektor_id = korektor.pns_id and korektor.deleted_at is null
WHERE 
    fds.deleted_at IS NULL
    AND fds.file_id = @id::varchar
    AND fds.ds_ok = true
ORDER BY fdc.korektor_ke ASC;

-- name: UpdateKorektorSuratKeputusanByID :exec
UPDATE 
    koreksi_surat_keputusan
SET 
    status_koreksi = COALESCE(@status_koreksi, status_koreksi),
    catatan_koreksi = COALESCE(@catatan_koreksi::text, catatan_koreksi)
WHERE 
    file_id = @id::varchar AND pegawai_korektor_id = @pns_id::varchar;

-- name: UpdateStatusSuratKeputusanByID :exec
UPDATE 
    surat_keputusan
SET 
    status_sk = COALESCE(@status_sk, status_sk),
    status_ttd = COALESCE(@status_ttd, status_ttd),
    status_kembali = COALESCE(@status_kembali, status_kembali),
    status_koreksi = COALESCE(@status_koreksi, status_koreksi),
    catatan = COALESCE(@catatan::text, catatan)
WHERE 
    file_id = @id::varchar;

-- name: InsertRiwayatSuratKeputusan :exec
INSERT INTO riwayat_surat_keputusan (
    file_id,
    nip_pemroses,
    tindakan,
    catatan_tindakan,
    akses_pengguna,
    waktu_tindakan,
    created_at,
    updated_at
) VALUES (
    @file_id::varchar,
    @nip_pemroses::varchar,
    @tindakan::varchar,
    @catatan_tindakan::text,
    @akses_pengguna::varchar,
    now(),
    now(),
    now()
);  

-- name: ListTandaTanganSuratKeputusanByPNSID :many
SELECT
    fds.file_id,
    p.nama as nama_pemilik_sk,
    p.nip_baru as nip_pemilik_sk,
    fds.kategori AS kategori_sk,
    fds.no_sk,
    fds.tanggal_sk,
    p.unor_id
FROM surat_keputusan fds
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fds.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip_pemilik')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip_pemilik')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('no_sk')::VARCHAR IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::VARCHAR || '%')
    AND fds.status_koreksi = 1 
    and (sqlc.narg('status_ttd')::integer is NULL or fds.status_ttd = sqlc.narg('status_ttd')::integer)
    AND fds.ttd_pegawai_id = @pns_id::varchar
    AND fds.ds_ok = true
ORDER BY fds.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountTandaTanganSuratKeputusanByPNSID :one
select
    count(1) as total
FROM surat_keputusan fds
JOIN pegawai p ON fds.nip_sk = p.nip_baru AND p.deleted_at IS NULL
LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
WHERE fds.deleted_at IS NULL
    AND (sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4)
    AND (sqlc.narg('nama_pemilik')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama_pemilik')::VARCHAR || '%')
    AND (sqlc.narg('nip_pemilik')::VARCHAR IS NULL OR fds.nip_sk = sqlc.narg('nip_pemilik')::VARCHAR)
    AND (sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER)
    AND (sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR)
    AND (sqlc.narg('kategori_sk')::VARCHAR is NULL OR fds.kategori ILIKE '%' || sqlc.narg('kategori_sk')::VARCHAR || '%')
    AND (sqlc.narg('no_sk')::VARCHAR IS NULL OR fds.no_sk ILIKE '%' || sqlc.narg('no_sk')::VARCHAR || '%')
    AND fds.status_koreksi = 1 
    and (sqlc.narg('status_ttd')::integer is NULL or fds.status_ttd = sqlc.narg('status_ttd')::integer)
    AND fds.ttd_pegawai_id = @pns_id::varchar
    AND fds.ds_ok = true;

-- name: ListTandaTanganSuratKeputusanAntreanByPNSID :many
SELECT 
    fds.kategori,
    fds.ttd_pegawai_id,
    COUNT(*) AS jumlah
FROM surat_keputusan fds
JOIN koreksi_surat_keputusan fdc  ON fds.file_id = fdc.file_id AND fdc.deleted_at IS NULL AND fdc.status_koreksi = 2
WHERE fds.status_ttd = 0
  AND fds.ds_ok = true
  AND fds.kategori NOT IN ('< Semua >', '< Pilih >')
  AND fds.ttd_pegawai_id = @pns_id::varchar
  AND fds.deleted_at IS NULL
  AND fds.ds_ok = true
GROUP BY 
    fds.kategori,
    fds.ttd_pegawai_id
LIMIT $1 OFFSET $2;

-- name: CountTandaTanganSuratKeputusanAntreanByPNSID :one
select
    count(distinct fds.kategori) as total
FROM surat_keputusan fds
JOIN koreksi_surat_keputusan fdc  ON fds.file_id = fdc.file_id AND fdc.deleted_at IS NULL AND fdc.status_koreksi = 2
WHERE fds.status_ttd = 0
  AND fds.ds_ok = true
  AND fds.kategori NOT IN ('< Semua >', '< Pilih >')
  AND fds.ttd_pegawai_id = @pns_id::varchar
  AND fds.deleted_at IS NULL
  AND fds.ds_ok = true;