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

-- name: ListPegawaiPPPK :many
SELECT 
	p.pns_id,
	p.nip_baru AS nip,
	p.nama,
	p.gelar_depan,
	p.gelar_belakang,
	p.foto,
	p.unor_id,
	p.status_cpns_pns,
	uk.nama_unor,
	ref_golongan_akhir.nama AS golongan,
	ref_jabatan.nama_jabatan AS jabatan,
	ref_kedudukan_hukum.nama as nama_kedudukuan_hukum
FROM pegawai p
	LEFT JOIN ref_unit_kerja as uk ON p.unor_id = uk.id
	JOIN ref_kedudukan_hukum
	    ON ref_kedudukan_hukum.id = p.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
	LEFT JOIN ref_jabatan
	    ON ref_jabatan.kode_jabatan = p.jabatan_instansi_id AND ref_jabatan.deleted_at IS NULL
	LEFT JOIN ref_golongan ref_golongan_akhir
	    ON ref_golongan_akhir.id = p.gol_id AND ref_golongan_akhir.deleted_at IS NULL
WHERE p.status_pegawai = 3
	AND (sqlc.narg('status_hukum')::varchar is null or ref_kedudukan_hukum.nama = sqlc.narg('status_hukum')::varchar)
	AND (
	    sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
	)
	AND (
	    sqlc.narg('keyword')::VARCHAR IS NULL
	    OR (
		p.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
		OR p.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
	    )
	)
	AND (
		sqlc.narg('nip')::VARCHAR IS NULL
		OR p.nip_baru = sqlc.narg('nip')::VARCHAR
	)
	AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER )
	AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
	AND (
		sqlc.narg('status_pns')::varchar[] IS NULL
		OR ( p.status_cpns_pns = ANY(sqlc.narg('status_pns')::VARCHAR[]) AND ref_kedudukan_hukum.nama <> @mpp::varchar )
	    )
	AND p.deleted_at IS NULL
ORDER BY p.nama ASC
LIMIT $1 OFFSET $2;

-- name: CountPegawaiPPPK :one
SELECT COUNT(1)
FROM pegawai p
	LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
	JOIN ref_kedudukan_hukum
	    ON ref_kedudukan_hukum.id = p.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
WHERE p.status_pegawai = 3
	AND (sqlc.narg('status_hukum')::varchar is null or ref_kedudukan_hukum.nama = sqlc.narg('status_hukum')::varchar)
	AND (
	    sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
	)
	AND (
	    sqlc.narg('keyword')::VARCHAR IS NULL
	    OR (
		p.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
		OR p.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
	    )
	)
	AND (
		sqlc.narg('nip')::VARCHAR IS NULL
		OR p.nip_baru = sqlc.narg('nip')::VARCHAR
	)
	AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER )
	AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
	AND (
		sqlc.narg('status_pns')::varchar[] IS NULL
		OR ( p.status_cpns_pns = ANY(sqlc.narg('status_pns')::VARCHAR[]) AND ref_kedudukan_hukum.nama <> @mpp::varchar )
	    )
	AND p.deleted_at IS NULL;

-- name: ListPegawaiNonAktif :many
SELECT 
	p.pns_id,
	p.nip_baru AS nip,
	p.nama,
	p.gelar_depan,
	p.gelar_belakang,
	p.foto,
	p.unor_id,
	p.status_cpns_pns,
	uk.nama_unor,
	ref_golongan_akhir.nama AS golongan,
	ref_jabatan.nama_jabatan AS jabatan,
	ref_kedudukan_hukum.nama as nama_kedudukuan_hukum
FROM pegawai p
	LEFT JOIN ref_unit_kerja as uk ON p.unor_id = uk.id
	JOIN ref_kedudukan_hukum
	    ON ref_kedudukan_hukum.id = p.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
	LEFT JOIN ref_jabatan
	    ON ref_jabatan.kode_jabatan = p.jabatan_instansi_id AND ref_jabatan.deleted_at IS NULL
	LEFT JOIN ref_golongan ref_golongan_akhir
	    ON ref_golongan_akhir.id = p.gol_id AND ref_golongan_akhir.deleted_at IS NULL
WHERE p.id is not null
  AND (p.kedudukan_hukum_id = '99' or p.status_pegawai = '3')
	AND (sqlc.narg('status_hukum')::varchar is null or ref_kedudukan_hukum.nama = sqlc.narg('status_hukum')::varchar)
	AND (
	    sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
	)
	AND (
	    sqlc.narg('keyword')::VARCHAR IS NULL
	    OR (
		p.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
		OR p.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
	    )
	)
	AND (
		sqlc.narg('nip')::VARCHAR IS NULL
		OR p.nip_baru = sqlc.narg('nip')::VARCHAR
	)
	AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER )
	AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
	AND (
		sqlc.narg('status_pns')::varchar[] IS NULL
		OR ( p.status_cpns_pns = ANY(sqlc.narg('status_pns')::VARCHAR[]) AND ref_kedudukan_hukum.nama <> @mpp::varchar )
	    )
	AND p.deleted_at IS NULL
ORDER BY p.nama ASC
LIMIT $1 OFFSET $2;

-- name: CountPegawaiNonAktif :one
SELECT COUNT(1)
FROM pegawai p
	LEFT JOIN ref_unit_kerja uk ON p.unor_id = uk.id AND uk.deleted_at IS NULL
	JOIN ref_kedudukan_hukum
	    ON ref_kedudukan_hukum.id = p.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
WHERE p.id is not null
  AND (p.kedudukan_hukum_id = '99' or p.status_pegawai = '3')
	AND (sqlc.narg('status_hukum')::varchar is null or ref_kedudukan_hukum.nama = sqlc.narg('status_hukum')::varchar)
	AND (
	    sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
	    OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
	)
	AND (
	    sqlc.narg('keyword')::VARCHAR IS NULL
	    OR (
		p.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
		OR p.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
	    )
	)
	AND (
		sqlc.narg('nip')::VARCHAR IS NULL
		OR p.nip_baru = sqlc.narg('nip')::VARCHAR
	)
	AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR p.gol_id = sqlc.narg('golongan_id')::INTEGER )
	AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR p.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
	AND (
		sqlc.narg('status_pns')::varchar[] IS NULL
		OR ( p.status_cpns_pns = ANY(sqlc.narg('status_pns')::VARCHAR[]) AND ref_kedudukan_hukum.nama <> @mpp::varchar )
	    )
	AND p.deleted_at IS NULL;


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

-- name: UpdateDataPegawai :exec
UPDATE pegawai
SET
    gelar_depan = @gelar_depan::varchar,
    nama = @nama::varchar,
    gelar_belakang = @gelar_belakang::varchar,
    nip_baru = @nip_baru::varchar,
    jenis_kelamin = @jenis_kelamin::varchar,
    nik = @nik::varchar,
    kk = @kk::varchar,
    tempat_lahir_id = @tempat_lahir_id::varchar,
    tanggal_lahir = @tanggal_lahir::date,
    tingkat_pendidikan_id = @tingkat_pendidikan_id::int2,
    pendidikan_id = @pendidikan_id::varchar,
    jenis_kawin_id = @jenis_kawin_id::int2,
    agama_id = @agama_id::int2,
    jenis_pegawai_id = @jenis_pegawai_id::int2,
    masa_kerja = @masa_kerja::varchar,
    jenis_jabatan_id = @jenis_jabatan_id::int2,
    jabatan_instansi_id = @jabatan_id::varchar,
    unor_id = @unor_id::varchar,
    lokasi_kerja_id = @lokasi_kerja_id::varchar,
    gol_awal_id = @gol_awal_id::int2,
    gol_id = @gol_id::int2,
    tmt_golongan = @tmt_golongan::date,
    tmt_pns = @tmt_pns::date,
    no_sk_cpns = @no_sk_cpns::varchar,
    status_cpns_pns = @status_cpns_pns::varchar,
    email_dikbud = @email_dikbud::varchar,
    email = @email::varchar,
    alamat = @alamat::varchar,
    no_hp = @no_hp::varchar,
    no_darurat = @no_darurat::varchar,
    no_surat_dokter = @no_surat_dokter::varchar,
    tanggal_surat_dokter = @tanggal_surat_dokter::date,
    no_bebas_narkoba = @no_bebas_narkoba::varchar,
    tanggal_bebas_narkoba = @tanggal_bebas_narkoba::date,
    no_catatan_polisi = @no_catatan_polisi::varchar,
    tanggal_catatan_polisi = @tanggal_catatan_polisi::date,
    akte_kelahiran = @akte_kelahiran::varchar,
    bpjs = @bpjs::varchar,
    npwp = @npwp::varchar,
    tanggal_npwp = @tanggal_npwp::date,
    no_taspen = @no_taspen::varchar,
    updated_at = now(),
    pns_id = @pns_id::varchar,
    mk_bulan = @mk_bulan::int2,
    mk_tahun = @mk_tahun::int2,
    mk_bulan_swasta = @mk_bulan_swasta::int2,
    mk_tahun_swasta = @mk_tahun_swasta::int2
WHERE nip_baru = @nip::varchar AND deleted_at IS NULL;

-- name: UpdateTTDPegawaiNIPByNIP :exec
UPDATE pegawai_ttd
SET 
    nip = @nip_baru::varchar,
    updated_at = now()
WHERE nip = @nip::varchar
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip)
);
    
