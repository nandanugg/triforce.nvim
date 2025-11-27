-- name: ListRiwayatAsesmenNineBox :many
select
  id,
  tahun,
  kesimpulan
from riwayat_asesmen_nine_box
where pns_nip = $1 and deleted_at is null
order by tahun desc nulls last
limit $2 offset $3;

-- name: CountRiwayatAsesmenNineBox :one
select count(1) from riwayat_asesmen_nine_box
where pns_nip = $1 and deleted_at is null;

-- name: UpdateRiwayatAsesmenNineBoxNamaNipByPNSID :exec
UPDATE riwayat_asesmen_nine_box
SET
    pns_nip = @nip_baru::varchar,
    nama = @nama::varchar,
    updated_at = now()
WHERE pns_nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM pns_nip)
    OR (@nama::varchar IS NOT NULL AND @nama::varchar IS DISTINCT FROM nama)
);
