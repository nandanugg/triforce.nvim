-- name: ListRefJenisKawin :many
select id, nama from ref_jenis_kawin
where deleted_at is null
limit $1 offset $2;

-- name: CountRefJenisKawin :one
select count(1) from ref_jenis_kawin
where deleted_at is null;
