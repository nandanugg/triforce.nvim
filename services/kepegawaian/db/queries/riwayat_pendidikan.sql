-- name: ListRiwayatPendidikan :many
select
    rp.id,
    rp.nama_sekolah,
    rp.tahun_lulus,
    rp.no_ijazah,
    rp.gelar_depan,
    rp.gelar_belakang,
    rp.tugas_belajar,
    rp.negara_sekolah,
    rp.tingkat_pendidikan_id,
    tk.nama as jenjang_pendidikan,
    rp.pendidikan_id,
    pend.nama as pendidikan
from riwayat_pendidikan rp
left join ref_tingkat_pendidikan tk on tk.id = rp.tingkat_pendidikan_id and tk.deleted_at is null
left join ref_pendidikan pend on pend.id = rp.pendidikan_id and pend.deleted_at is null
where rp.nip = $1 and rp.deleted_at is null
order by rp.tahun_lulus desc nulls last
limit $2 offset $3;

-- name: GetRiwayatPendidikan :one
select
    rp.id,
    rp.nama_sekolah,
    rp.tahun_lulus,
    rp.no_ijazah,
    rp.gelar_depan,
    rp.gelar_belakang,
    rp.tugas_belajar,
    rp.negara_sekolah,
    rp.tingkat_pendidikan_id,
    tk.nama as tingkat_pendidikan,
    rp.pendidikan_id,
    pend.nama as pendidikan
from riwayat_pendidikan rp
left join ref_tingkat_pendidikan tk on tk.id = rp.tingkat_pendidikan_id and tk.deleted_at is null
left join ref_pendidikan pend on pend.id = rp.pendidikan_id and pend.deleted_at is null
where rp.nip = @nip::varchar and rp.id = @id and rp.deleted_at is null;

-- name: CountRiwayatPendidikan :one
select count(1)
from riwayat_pendidikan
where nip = $1 and deleted_at is null;

-- name: GetBerkasRiwayatPendidikan :one
select file_base64 from riwayat_pendidikan
where nip = $1 and id = $2 and deleted_at is null;

-- name: CreateRiwayatPendidikan :one
insert into riwayat_pendidikan
    (tingkat_pendidikan_id, pendidikan_id, nama_sekolah, tahun_lulus, no_ijazah, gelar_depan, gelar_belakang, negara_sekolah, tugas_belajar, pns_id, nip) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning id;

-- name: UpdateRiwayatPendidikan :execrows
update riwayat_pendidikan
set
    tingkat_pendidikan_id = $1,
    pendidikan_id = $2,
    nama_sekolah = $3,
    tahun_lulus = $4,
    no_ijazah = $5,
    gelar_depan = $6,
    gelar_belakang = $7,
    negara_sekolah = $8,
    tugas_belajar = $9,
    updated_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatPendidikan :execrows
update riwayat_pendidikan
set deleted_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;

-- name: UploadBerkasRiwayatPendidikan :execrows
update riwayat_pendidikan
set
    file_base64 = $1,
    updated_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;
