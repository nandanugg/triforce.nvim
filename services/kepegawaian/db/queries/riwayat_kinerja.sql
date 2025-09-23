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
