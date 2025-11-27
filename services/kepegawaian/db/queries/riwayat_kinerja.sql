-- name: ListRiwayatKinerja :many
select
  id,
  tahun,
  rating_hasil_kerja,
  rating_perilaku_kerja,
  predikat_kinerja
from riwayat_kinerja
where nip = $1 and deleted_at is null
order by tahun desc nulls last
limit $2 offset $3;

-- name: CountRiwayatKinerja :one
select count(1) from riwayat_kinerja
where nip = $1 and deleted_at is null;

-- name: UpdateRiwayatKinerjaNamaNipByPNSID :exec
UPDATE riwayat_kinerja
SET     
    nip = @nip_baru::varchar,
    nama = @nama::varchar,
    updated_at = now()
WHERE nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip)
    OR (@nama::varchar IS NOT NULL AND @nama::varchar IS DISTINCT FROM nama)
);