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
