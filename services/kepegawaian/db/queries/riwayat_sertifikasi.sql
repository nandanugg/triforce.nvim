-- name: ListRiwayatSertifikasi :many
select id, nama_sertifikasi, tahun, deskripsi from riwayat_sertifikasi
where nip = $1 and deleted_at is null
order by tahun desc nulls last
limit $2 offset $3;

-- name: CountRiwayatSertifikasi :one
select count(1) from riwayat_sertifikasi
where nip = $1 and deleted_at is null;

-- name: GetBerkasRiwayatSertifikasi :one
select file_base64 from riwayat_sertifikasi
where nip = $1 and id = $2 and deleted_at is null;

-- name: CreateRiwayatSertifikasi :one
insert into 
    riwayat_sertifikasi (nip, tahun, nama_sertifikasi, deskripsi) 
values 
    ($1, $2, $3, $4)
returning id;

-- name: UpdateBerkasRiwayatSertifikasiByIDAndNIP :execrows
update riwayat_sertifikasi
set
    file_base64 = $1,
    updated_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;

-- name: UpdateRiwayatSertifikasiByIDAndNIP :execrows
update riwayat_sertifikasi
set
    tahun = $1,
    nama_sertifikasi = $2,
    deskripsi = $3,
    updated_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatSertifikasiByIDAndNIP :execrows
update riwayat_sertifikasi
set
    deleted_at = now()
where id = @id and nip = @nip::varchar and deleted_at is null;

-- name: UpdateRiwayatSertifikasiNamaNipByNIP :exec
UPDATE riwayat_sertifikasi
SET     
    nip = @nip_baru::varchar,
    updated_at = now()
WHERE nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip)
);