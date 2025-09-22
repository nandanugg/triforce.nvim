-- name: ListRefJenisPegawai :many
select id, nama from ref_jenis_pegawai
where deleted_at is null
limit $1 offset $2;

-- name: CountRefJenisPegawai :one
select count(1) from ref_jenis_pegawai
where deleted_at is null;
