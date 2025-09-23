-- name: ListRiwayatHukumanDisiplin :many
SELECT
  rh.id,
  COALESCE(rh.nama_jenis_hukuman, rjh.nama) as jenis_hukuman,
  rg.nama as nama_golongan,
  rg.nama_pangkat,
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
