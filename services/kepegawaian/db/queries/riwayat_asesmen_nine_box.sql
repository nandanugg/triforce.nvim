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
