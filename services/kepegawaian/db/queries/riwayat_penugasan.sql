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
