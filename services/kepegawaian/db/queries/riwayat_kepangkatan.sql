-- name: ListRiwayatKepangkatan :many
select
    rg.id,
    ref_jenis_kp.id as jenis_kp_id,
    ref_jenis_kp.nama as nama_jenis_kp,
    ref_golongan.id as golongan_id,
    ref_golongan.nama as nama_golongan,
    ref_golongan.nama_pangkat as nama_golongan_pangkat,
    rg.tmt_golongan,
    rg.sk_nomor,
    rg.sk_tanggal,
    rg.mk_golongan_tahun,
    rg.mk_golongan_bulan,
    rg.no_bkn,
    rg.tanggal_bkn,
    rg.jumlah_angka_kredit_tambahan,
    rg.jumlah_angka_kredit_utama
from riwayat_golongan rg
LEFT join ref_jenis_kp on rg.jenis_kp_id = ref_jenis_kp.id and ref_jenis_kp.deleted_at is null
LEFT join ref_golongan on rg.golongan_id = ref_golongan.id and ref_golongan.deleted_at is null
where rg.deleted_at is null and rg.pns_nip = @pns_nip::varchar
order by rg.tmt_golongan desc
limit $1 offset $2;

-- name: CountRiwayatKepangkatan :one
select count(*) from riwayat_golongan rg
where rg.deleted_at is null and rg.pns_nip = @pns_nip::varchar;

-- name: GetBerkasRiwayatKepangkatan :one
select file_base64 from riwayat_golongan
where pns_nip = $1 and id = $2 and deleted_at is null;
