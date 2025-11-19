-- name: ListRiwayatPenugasan :many
select
  id,
  tipe_jabatan,
  nama_jabatan,
  deskripsi_jabatan,
  tanggal_mulai,
  tanggal_selesai,
  is_menjabat
from riwayat_penugasan
where nip = $1 and deleted_at is null
order by tanggal_mulai desc nulls last
limit $2 offset $3;

-- name: CountRiwayatPenugasan :one
select count(1) from riwayat_penugasan
where nip = $1 and deleted_at is null;

-- name: GetBerkasRiwayatPenugasan :one
select file_base64 from riwayat_penugasan
where nip = $1 and id = $2 and deleted_at is null;

-- name: CreateRiwayatPenugasan :one
insert into riwayat_penugasan (tipe_jabatan, nama_jabatan, deskripsi_jabatan, tanggal_mulai, tanggal_selesai, nip, is_menjabat)
values ($1, $2, $3, $4, $5, $6, $7)
returning id;

-- name: UpdateRiwayatPenugasan :execrows
update 
  riwayat_penugasan
set 
  tipe_jabatan = $1, 
  nama_jabatan = $2, 
  deskripsi_jabatan = $3, 
  tanggal_mulai = $4, 
  tanggal_selesai = $5, 
  is_menjabat = $6,
  updated_at = now()
where 
  id = @id and nip = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatPenugasan :execrows
update riwayat_penugasan
set deleted_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;

-- name: UploadBerkasRiwayatPenugasan :execrows
update riwayat_penugasan
set file_base64 = $1, updated_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;