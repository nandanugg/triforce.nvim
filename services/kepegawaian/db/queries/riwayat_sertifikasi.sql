-- name: ListRiwayatSertifikasi :many
select id, nama_sertifikasi, tahun, deskripsi from riwayat_sertifikasi
where nip = $1 and deleted_at is null
order by tahun desc
limit $2 offset $3;

-- name: CountRiwayatSertifikasi :one
select count(1) from riwayat_sertifikasi
where nip = $1 and deleted_at is null;

-- name: GetBerkasRiwayatSertifikasi :one
select file_base64 from riwayat_sertifikasi
where nip = $1 and id = $2 and deleted_at is null;
