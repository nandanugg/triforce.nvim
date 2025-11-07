-- name: ListRiwayatKepangkatan :many
select
    rg.id,
    ref_jenis_kenaikan_pangkat.id as jenis_kp_id,
    ref_jenis_kenaikan_pangkat.nama as nama_jenis_kp,
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
LEFT join ref_jenis_kenaikan_pangkat on rg.jenis_kp_id = ref_jenis_kenaikan_pangkat.id and ref_jenis_kenaikan_pangkat.deleted_at is null
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

-- name: CreateRiwayatKepangkatan :one
insert into riwayat_golongan
    (id, jenis_kp_id, kode_jenis_kp, jenis_kp, golongan_id, golongan_nama, pangkat_nama, tmt_golongan, sk_nomor, sk_tanggal, no_bkn, tanggal_bkn, jumlah_angka_kredit_utama, jumlah_angka_kredit_tambahan, mk_golongan_tahun, mk_golongan_bulan, pns_id, pns_nip, pns_nama) values
    (uuid_generate_v4(), sqlc.narg('jenis_kp_id')::int, sqlc.narg('jenis_kp_id')::varchar, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
returning id;

-- name: UpdateRiwayatKepangkatan :execrows
update riwayat_golongan
set
    jenis_kp_id = sqlc.narg('jenis_kp_id')::int,
    kode_jenis_kp = sqlc.narg('jenis_kp_id')::varchar,
    jenis_kp = $1,
    golongan_id = $2,
    golongan_nama = $3,
    pangkat_nama = $4,
    tmt_golongan = $5,
    sk_nomor = $6,
    sk_tanggal = $7,
    no_bkn = $8,
    tanggal_bkn = $9,
    jumlah_angka_kredit_utama = $10,
    jumlah_angka_kredit_tambahan = $11,
    mk_golongan_tahun = $12,
    mk_golongan_bulan = $13,
    updated_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatKepangkatan :execrows
update riwayat_golongan
set deleted_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: UploadBerkasRiwayatKepangkatan :execrows
update riwayat_golongan
set
    file_base64 = $1,
    updated_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;
