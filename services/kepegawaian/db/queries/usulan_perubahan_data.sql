-- name: ListUnreadUsulanPerubahanDataByNIP :many
select
  id,
  jenis_data,
  data_id,
  perubahan_data,
  action,
  status,
  catatan
from usulan_perubahan_data
where nip = $1 and jenis_data = $2
  and read_at is null and deleted_at is null
order by updated_at desc
limit $3 offset $4;

-- name: CountUnreadUsulanPerubahanDataByNIP :one
select count(1)
from usulan_perubahan_data
where nip = $1 and jenis_data = $2
  and read_at is null and deleted_at is null;

-- name: ListPendingUsulanPerubahanData :many
select
  upd.id,
  upd.nip,
  upd.jenis_data,
  upd.created_at,
  p.gelar_depan as gelar_depan_pegawai,
  p.gelar_belakang as gelar_belakang_pegawai,
  p.nama as nama_pegawai,
  p.unor_id as unor_id_pegawai
from usulan_perubahan_data upd
join pegawai p on p.nip_baru = upd.nip and p.deleted_at is null
where upd.status = 'Diusulkan' and upd.deleted_at is null
  and (@nama::varchar = '' or p.nama ilike '%' || @nama::varchar || '%')
  and (@nip::varchar = '' or upd.nip like @nip::varchar || '%')
  and (@jenis_data::varchar = '' or upd.jenis_data = @jenis_data::varchar)
  and (@kode_jabatan::varchar = '' or p.jabatan_instansi_id = @kode_jabatan::varchar)
  and (sqlc.narg('golongan_id')::int2 is null or p.gol_id = sqlc.narg('golongan_id')::int2)
  and (@unit_kerja_id::varchar = '' or p.unor_id = @unit_kerja_id::varchar)
order by upd.created_at asc
limit $1 offset $2;

-- name: CountPendingUsulanPerubahanData :one
select count(1)
from usulan_perubahan_data upd
join pegawai p on p.nip_baru = upd.nip and p.deleted_at is null
where upd.status = 'Diusulkan' and upd.deleted_at is null
  and (@nama::varchar = '' or p.nama ilike '%' || @nama::varchar || '%')
  and (@nip::varchar = '' or upd.nip like @nip::varchar || '%')
  and (@jenis_data::varchar = '' or upd.jenis_data = @jenis_data::varchar)
  and (@kode_jabatan::varchar = '' or p.jabatan_instansi_id = @kode_jabatan::varchar)
  and (sqlc.narg('golongan_id')::int2 is null or p.gol_id = sqlc.narg('golongan_id')::int2)
  and (@unit_kerja_id::varchar = '' or p.unor_id = @unit_kerja_id::varchar);

-- name: GetUsulanPerubahanData :one
select
  upd.id,
  upd.nip,
  upd.jenis_data,
  upd.data_id,
  upd.perubahan_data,
  upd.action,
  upd.status,
  upd.catatan,
  upd.created_at,
  p.gelar_depan as gelar_depan_pegawai,
  p.gelar_belakang as gelar_belakang_pegawai,
  p.nama as nama_pegawai,
  case when kh.is_pppk then g.gol_pppk else g.nama end as golongan_pegawai,
  j.nama_jabatan as jabatan_pegawai,
  p.foto as foto_pegawai,
  case
    when p.status_cpns_pns = 'P' then 'PNS'
    when p.status_cpns_pns = 'C' then 'CPNS'
    else p.status_cpns_pns
  end as status_pns_pegawai,
  p.unor_id as unor_id_pegawai
from usulan_perubahan_data upd
join pegawai p on p.nip_baru = upd.nip and p.deleted_at is null
left join ref_jabatan j on j.kode_jabatan = p.jabatan_instansi_id and j.deleted_at is null
left join ref_golongan g on g.id = p.gol_id and g.deleted_at is null
left join ref_kedudukan_hukum kh on kh.id = p.kedudukan_hukum_id and kh.deleted_at is null
where upd.id = $1 and upd.jenis_data = $2 and upd.deleted_at is null;

-- name: UpdateStatusUsulanPerubahanData :exec
update usulan_perubahan_data
set
  status = $3,
  catatan = $4,
  updated_at = now()
where id = $1 and jenis_data = $2 and deleted_at is null;

-- name: CreateUsulanPerubahanData :one
insert into usulan_perubahan_data
  (nip, jenis_data, data_id, perubahan_data, action) values
  ($1, $2, $3, $4, $5)
returning id;

-- name: MarkAsReadUsulanPerubahanData :exec
update usulan_perubahan_data
set read_at = now()
where id = $1 and jenis_data = $2 and deleted_at is null;

-- name: DeleteUsulanPerubahanData :exec
update usulan_perubahan_data
set deleted_at = now()
where id = $1 and jenis_data = $2 and deleted_at is null;
