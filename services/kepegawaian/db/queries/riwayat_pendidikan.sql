-- name: ListRiwayatPendidikan :many
select
    rp.id,
    rp.nama_sekolah,
    rp.tahun_lulus,
    rp.no_ijazah,
    rp.gelar_depan,
    rp.gelar_belakang,
    rp.tugas_belajar,
    rp.negara_sekolah,
    tk.nama as jenjang_pendidikan,
    pend.nama as pendidikan
from riwayat_pendidikan rp
left join ref_tingkat_pendidikan tk on tk.id = rp.tingkat_pendidikan_id and tk.deleted_at is null
left join ref_pendidikan pend on pend.id = rp.pendidikan_id and pend.deleted_at is null
where rp.nip = $1 and rp.deleted_at is null
order by rp.tahun_lulus desc
limit $2 offset $3;

-- name: CountRiwayatPendidikan :one
select count(1)
from riwayat_pendidikan
where nip = $1 and deleted_at is null;

-- name: GetBerkasRiwayatPendidikan :one
select file_base64 from riwayat_pendidikan
where nip = $1 and id = $2 and deleted_at is null;
