-- name: GetProfilePegawaiByPNSID :one
select
  p.nip_lama,
  p.nip_baru,
  p.gelar_depan,
  p.gelar_belakang,
  p.nama,
  p.unor_id,
  rj.nama_jabatan as jabatan,
  rg.nama_pangkat as pangkat,
  case when rkh.is_pppk then rg.gol_pppk else rg.nama end as golongan,
  p.foto
from pegawai p
left join ref_jabatan rj on rj.kode_jabatan = p.jabatan_instansi_id and rj.deleted_at is null
left join ref_golongan rg on rg.id = p.gol_id and rg.deleted_at is null
left join ref_kedudukan_hukum rkh on rkh.id = p.kedudukan_hukum_id and rkh.deleted_at is null
where pns_id = $1 and p.deleted_at is null;

-- name: ListPegawaiAktif :many
SELECT
    pegawai.pns_id,
    pegawai.nip_baru AS nip,
    pegawai.nama,
    pegawai.gelar_depan,
    pegawai.gelar_belakang,
    ref_golongan_akhir.nama AS golongan,
    ref_jabatan.nama_jabatan AS jabatan,
    pegawai.unor_id,
    ref_kedudukan_hukum.nama as nama_kedudukuan_hukum,
    status_cpns_pns,
    pegawai.foto
FROM pegawai
LEFT JOIN ref_unit_kerja uk
    ON pegawai.unor_id = uk.id AND uk.deleted_at IS NULL
JOIN ref_kedudukan_hukum
    ON ref_kedudukan_hukum.id = pegawai.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
LEFT JOIN ref_jabatan
    ON ref_jabatan.kode_jabatan = pegawai.jabatan_instansi_id AND ref_jabatan.deleted_at IS NULL
LEFT JOIN ref_golongan ref_golongan_akhir
    ON ref_golongan_akhir.id = pegawai.gol_id AND ref_golongan_akhir.deleted_at IS NULL
WHERE
    ref_kedudukan_hukum.is_pegawai_aktif = TRUE
    AND (sqlc.narg('status_hukum')::varchar is null or ref_kedudukan_hukum.nama = sqlc.narg('status_hukum')::varchar)
    AND (
        sqlc.narg('keyword')::VARCHAR IS NULL
        OR (
            pegawai.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
            OR pegawai.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
        )
    )
    AND (
        sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
    )
    AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR pegawai.gol_id = sqlc.narg('golongan_id')::INTEGER )
    AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR pegawai.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
    AND (
        sqlc.narg('status_pns')::varchar[] IS NULL
        OR ( pegawai.status_cpns_pns = ANY(sqlc.narg('status_pns')::VARCHAR[]) AND ref_kedudukan_hukum.nama <> @mpp::varchar )
    )
    AND pegawai.deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: ListPegawaiPPNPN :many
SELECT p.*, ruk.nama_unor
FROM pegawai p
         LEFT JOIN ref_unit_kerja as ruk ON p.unor_id = ruk.id
WHERE p.status_pegawai = 3
	AND sqlc.narg('unit_kerja_id')::VARCHAR IS NULL OR ruk.id = sqlc.narg('unit_kerja_id')::VARCHAR 
	AND sqlc.narg('nama')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama')::VARCHAR || '%'
	AND sqlc.narg('nip')::VARCHAR IS NULL OR p.nip_baru = sqlc.narg('nip')::VARCHAR
ORDER BY p.nama ASC
LIMIT 10;

-- name: ListPegawaiNonAktif :many
SELECT p.id,
       p.pns_id,
       p.nip_baru,
       p.nama,
       ruk.nama_unor,
       ref_golongan.nama as nama_golongan,
       nama_pangkat,
       eselon_4,
       eselon_3,
       eselon_2,
       eselon_1,
       kategori_jabatan
FROM pegawai p
         LEFT JOIN ref_unit_kerja as ruk ON p.unor_id = ruk.id
         LEFT JOIN ref_golongan ON p.gol_id = ref_golongan.id
         LEFT JOIN ref_jabatan ON p.jabatan_instansi_id = ref_jabatan.kode_jabatan
WHERE p.id is not null
  AND (p.kedudukan_hukum_id = '99' or p.status_pegawai = '3')
	AND sqlc.narg('unit_kerja_id')::VARCHAR IS NULL OR ruk.id = sqlc.narg('unit_kerja_id')::VARCHAR 
	AND sqlc.narg('nama')::VARCHAR IS NULL OR p.nama ILIKE '%' || sqlc.narg('nama')::VARCHAR || '%'
	AND sqlc.narg('nip')::VARCHAR IS NULL OR p.nip_baru = sqlc.narg('nip')::VARCHAR
	AND sqlc.narg('golongan_id')::VARCHAR IS NULL OR p.gol_id = sqlc.narg('golongan_id')::VARCHAR
	AND sqlc.narg('tingkat_pendidikan_id')::VARCHAR IS NULL OR p.tingkat_pendidikan_id = sqlc.narg('tingkat_pendidikan_id')::VARCHAR
	AND sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_id = sqlc.narg('jabatan_id')::VARCHAR
ORDER BY p.nama ASC
LIMIT 10;

-- name: CountPegawaiAktif :one
SELECT
  COUNT(1)
FROM pegawai
LEFT JOIN ref_unit_kerja uk
    ON pegawai.unor_id = uk.id AND uk.deleted_at IS NULL
JOIN ref_kedudukan_hukum
    ON ref_kedudukan_hukum.id = pegawai.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
WHERE
    ref_kedudukan_hukum.is_pegawai_aktif = TRUE
    AND (sqlc.narg('status_hukum')::varchar is null or ref_kedudukan_hukum.nama = sqlc.narg('status_hukum')::varchar)
    AND (
        sqlc.narg('keyword')::VARCHAR IS NULL
        OR (
            pegawai.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
            OR pegawai.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
        )
    )
    AND (
        sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
        OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
    )
    AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR pegawai.gol_id = sqlc.narg('golongan_id')::INTEGER )
    AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR pegawai.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
    AND (
        sqlc.narg('status_pns')::varchar[] IS NULL
        OR ( pegawai.status_cpns_pns = ANY(sqlc.narg('status_pns')::VARCHAR[]) AND ref_kedudukan_hukum.nama <> @mpp::varchar )
    )
    AND pegawai.deleted_at IS NULL;

-- name: GetPegawaiPNSIDByNIP :one
SELECT pns_id FROM pegawai WHERE nip_baru = @nip::varchar AND deleted_at IS NULL;

-- name: GetPegawaiTTDByNIP :one
SELECT base64ttd FROM pegawai_ttd WHERE nip = @nip::varchar;

-- name: GetPegawaiNIKByNIP :one
SELECT nik::text FROM pegawai WHERE nip_baru = @nip::varchar AND deleted_at IS NULL;

-- name: GetPegawaiByNIP :one
select
    pegawai.id as id,
    pns_id,
    nip_baru as nip,
    pegawai.nama as nama,
    tanggal_lahir,
    coalesce(ref_lokasi.nama, tempat_lahir) as tempat_lahir
FROM pegawai
LEFT JOIN ref_lokasi on ref_lokasi.id = pegawai.tempat_lahir_id and ref_lokasi.deleted_at is null
where nip_baru = @nip::varchar and pegawai.deleted_at is null;
