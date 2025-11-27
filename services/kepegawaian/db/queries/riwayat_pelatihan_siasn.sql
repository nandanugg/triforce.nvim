-- name: ListRiwayatPelatihanSIASN :many
SELECT
    rd.id,
    rd.jenis_diklat_id,
    rjd.jenis_diklat,
    rd.nama_diklat,
    rd.no_sertifikat,
    rd.tanggal_mulai,
    rd.tanggal_selesai,
    rd.tahun_diklat,
    rd.durasi_jam,
    rd.institusi_penyelenggara
FROM riwayat_diklat rd
LEFT JOIN ref_jenis_diklat rjd ON rd.jenis_diklat_id = rjd.id AND rjd.deleted_at IS NULL
WHERE rd.nip_baru = $1 AND rd.deleted_at IS NULL
ORDER BY rd.tanggal_selesai DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountRiwayatPelatihanSIASN :one
SELECT count(*) FROM riwayat_diklat
WHERE nip_baru = $1 AND deleted_at IS NULL;

-- name: GetBerkasRiwayatPelatihanSIASN :one
select file_base64 from riwayat_diklat
where nip_baru = $1 and id = $2 and deleted_at is null;

-- name: CreateRiwayatPelatihanSIASN :one
insert into riwayat_diklat
    (nama_diklat, jenis_diklat_id, jenis_diklat, institusi_penyelenggara, no_sertifikat, tanggal_mulai, tanggal_selesai, tahun_diklat, durasi_jam, sudah_kirim_siasn, pns_orang_id, nip_baru) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, null, $10, $11)
returning id;

-- name: UpdateRiwayatPelatihanSIASN :execrows
update riwayat_diklat
set
    nama_diklat = $1,
    jenis_diklat_id = $2,
    jenis_diklat = $3,
    institusi_penyelenggara = $4,
    no_sertifikat = $5,
    tanggal_mulai = $6,
    tanggal_selesai = $7,
    tahun_diklat = $8,
    durasi_jam = $9,
    updated_at = now()
where id = @id and nip_baru = @nip::varchar and deleted_at is null;

-- name: DeleteRiwayatPelatihanSIASN :execrows
update riwayat_diklat
set deleted_at = now()
where id = @id and nip_baru = @nip::varchar and deleted_at is null;

-- name: UploadBerkasRiwayatPelatihanSIASN :execrows
update riwayat_diklat
set
    file_base64 = $1,
    updated_at = now()
where id = @id and nip_baru = @nip::varchar and deleted_at is null;

-- name: UpdateRiwayatPelatihanSIASNNamaNipByPNSID :exec
UPDATE riwayat_diklat
SET     
    nip_baru = @nip_baru::varchar,
    updated_at = now()
WHERE pns_orang_id = @pns_id::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip_baru)
);
