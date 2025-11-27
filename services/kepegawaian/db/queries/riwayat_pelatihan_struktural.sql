-- name: ListRiwayatPelatihanStruktural :many
SELECT
    rs.id,
    rs.nama_diklat,
    rs.tanggal,
    rs.tahun,
    rs.nomor,
    rs.lama
FROM riwayat_diklat_struktural rs
WHERE rs.pns_nip = $1
AND rs.deleted_at IS NULL
ORDER BY rs."tahun" DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanStruktural :one
SELECT COUNT(1) AS total
FROM riwayat_diklat_struktural rs
WHERE rs.pns_nip = $1
  AND rs.deleted_at IS NULL;

-- name: GetBerkasRiwayatPelatihanStruktural :one
select file_base64 from riwayat_diklat_struktural
where pns_nip = $1 and id = $2 and deleted_at is null;

-- name: CreateRiwayatPelatihanStruktural :one
insert into riwayat_diklat_struktural
    (id, nama_diklat, tanggal, tahun, lama, nomor, pns_id, pns_nip, pns_nama) values
    (public.uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8)
returning id;

-- name: UpdateRiwayatPelatihanStruktural :execrows
update riwayat_diklat_struktural
set
    nama_diklat = $1,
    tanggal = $2,
    tahun = $3,
    lama = $4,
    nomor = $5,
    updated_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatPelatihanStruktural :execrows
update riwayat_diklat_struktural
set deleted_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: UploadBerkasRiwayatPelatihanStruktural :execrows
update riwayat_diklat_struktural
set
    file_base64 = $1,
    updated_at = now()
where id = @id and pns_nip = @nip::varchar and deleted_at is null;

-- name: UpdateRiwayatPelatihanStrukturalNamaNipByPNSID :exec
UPDATE riwayat_diklat_struktural
SET     
    pns_nip = @nip_baru::varchar,
    pns_nama = @nama::varchar,
    updated_at = now()
WHERE pns_id = @pns_id::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM pns_nip)
    OR (@nama::varchar IS NOT NULL AND @nama::varchar IS DISTINCT FROM pns_nama)
);
