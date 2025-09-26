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
  case when rkh.is_pppk then rg.gol_pppk else rg.nama end as golongan
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
    CASE WHEN 
        ref_kedudukan_hukum.nama = @mpp::varchar then 'MPP'
    ELSE 
        pegawai.status_cpns_pns
    END::text AS status
FROM pegawai
LEFT JOIN unit_kerja uk 
    ON pegawai.unor_id = uk.id AND uk.deleted_at IS NULL
LEFT JOIN ref_kedudukan_hukum
    ON ref_kedudukan_hukum.id = pegawai.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
LEFT JOIN ref_jabatan
    ON ref_jabatan.kode_jabatan = pegawai.jabatan_instansi_id AND ref_jabatan.deleted_at IS NULL
LEFT JOIN ref_golongan ref_golongan_akhir
    ON ref_golongan_akhir.id = pegawai.gol_id AND ref_golongan_akhir.deleted_at IS NULL
WHERE 
    ref_kedudukan_hukum.nama = ANY(@status_hukum::varchar[])
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
    AND ( @status_hukum::varchar[] = ARRAY[@mpp::varchar] OR sqlc.narg('status_pns')::varchar IS NULL OR pegawai.status_cpns_pns = sqlc.narg('status_pns')::VARCHAR )
    AND pegawai.deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountPegawaiAktif :one
SELECT 
  COUNT(1)
FROM pegawai
LEFT JOIN unit_kerja uk 
    ON pegawai.unor_id = uk.id AND uk.deleted_at IS NULL
LEFT JOIN ref_kedudukan_hukum
    ON ref_kedudukan_hukum.id = pegawai.kedudukan_hukum_id AND ref_kedudukan_hukum.deleted_at IS NULL
WHERE 
    ref_kedudukan_hukum.nama = ANY(@status_hukum::varchar[])
    AND (
        sqlc.narg('keyword')::VARCHAR IS NULL
        OR (
            pegawai.nama ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
            OR pegawai.nip_baru ILIKE '%' || sqlc.narg('keyword')::VARCHAR || '%'
        )
    )
    AND (
        sqlc.narg('unit_kerja_id')::VARCHAR IS NULL
        OR (
            uk.id is not null AND (
                sqlc.narg('unit_kerja_id')::VARCHAR = uk.id
                OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_1
                OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_2
                OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_3
                OR sqlc.narg('unit_kerja_id')::VARCHAR = uk.eselon_4
            )
        )
    )
    AND ( sqlc.narg('golongan_id')::INTEGER IS NULL OR pegawai.gol_id = sqlc.narg('golongan_id')::INTEGER )
    AND ( sqlc.narg('jabatan_id')::VARCHAR IS NULL OR pegawai.jabatan_instansi_id = sqlc.narg('jabatan_id')::VARCHAR )
    AND ( @status_hukum::varchar[] = ARRAY[@mpp::varchar] OR sqlc.narg('status_pns')::varchar IS NULL OR pegawai.status_cpns_pns = sqlc.narg('status_pns')::VARCHAR )
    AND pegawai.deleted_at IS NULL;