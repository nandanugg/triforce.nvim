-- name: ListRiwayatHukumanDisiplin :many
SELECT
  rh.id,
  COALESCE(rh.nama_jenis_hukuman, rjh.nama) as jenis_hukuman,
  rh.jenis_hukuman_id,
  rg.nama as nama_golongan,
  rg.nama_pangkat,
  rh.golongan_id,
  tanggal_mulai_hukuman,
  tanggal_akhir_hukuman,
  masa_tahun,
  masa_bulan,
  sk_nomor,
  sk_tanggal,
  no_pp,
  no_sk_pembatalan,
  tanggal_sk_pembatalan
FROM riwayat_hukuman_disiplin rh
LEFT JOIN ref_jenis_hukuman rjh ON rh.jenis_hukuman_id=rjh.id AND rjh.deleted_at is null
LEFT JOIN ref_golongan rg ON rh.golongan_id=rg.id AND rg.deleted_at is null
WHERE pns_nip = $1
  AND rh.deleted_at is null
ORDER BY sk_tanggal DESC
LIMIT $2 OFFSET $3;

-- name: GetBerkasRiwayatHukumanDisiplin :one
SELECT file_base64
FROM riwayat_hukuman_disiplin rh
WHERE pns_nip = $1
  AND rh.id = $2
  AND rh.deleted_at is null;

-- name: CountRiwayatHukumanDisiplin :one
SELECT COUNT(1)
FROM riwayat_hukuman_disiplin rh
WHERE pns_nip = $1
  AND rh.deleted_at is null;

-- name: CreateRiwayatHukumanDisiplin :one
INSERT INTO riwayat_hukuman_disiplin (
  pns_id,
  pns_nip,
  nama,
  golongan_id,
  nama_golongan,
  jenis_hukuman_id,
  nama_jenis_hukuman,
  sk_nomor,
  sk_tanggal,
  tanggal_mulai_hukuman,
  masa_tahun,
  masa_bulan,
  tanggal_akhir_hukuman,
  no_pp,
  no_sk_pembatalan,
  tanggal_sk_pembatalan
)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  $15,
  $16
)
RETURNING id;

-- name: UpdateRiwayatHukumanDisiplin :execrows
UPDATE riwayat_hukuman_disiplin
SET
  golongan_id = $1,
  jenis_hukuman_id = $2,
  nama_golongan = $3,
  nama_jenis_hukuman = $4,
  sk_nomor = $5,
  sk_tanggal = $6,
  tanggal_mulai_hukuman = $7,
  masa_tahun = $8,
  masa_bulan = $9,
  tanggal_akhir_hukuman = $10,
  no_pp = $11,
  no_sk_pembatalan = $12,
  tanggal_sk_pembatalan = $13,
  updated_at = now()
WHERE id = @id::integer and pns_nip = @nip::varchar
  AND deleted_at is null;

-- name: DeleteRiwayatHukumanDisiplin :execrows
UPDATE riwayat_hukuman_disiplin
SET
  deleted_at = now()
WHERE id = @id::integer and pns_nip = @nip::varchar
  AND deleted_at is null;

-- name: UploadBerkasRiwayatHukumanDisiplin :execrows
UPDATE riwayat_hukuman_disiplin
SET
  file_base64 = $1,
  updated_at = now()
WHERE id = @id and pns_nip = @nip::varchar
  AND deleted_at is null;