-- name: GetRefJabatan :many
select id, kode_jabatan, nama_jabatan from ref_jabatan
LIMIT $1 OFFSET $2;