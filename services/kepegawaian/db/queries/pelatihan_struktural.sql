-- name: ListPelatihanStruktural :many
SELECT 
	rjd.jenis_diklat,
	rd.id,
	rd.nama_diklat,
	rd.institusi_penyelenggara,
	rd.no_sertifikat,
	rd.tanggal_mulai,
	rd.tanggal_selesai,
	rd.tahun_diklat,
	rd.durasi_jam
FROM 
	riwayat_diklat rd
JOIN 
	ref_jenis_diklat rjd ON rd.jenis_diklat_id = rjd.id
WHERE 
	rd.deleted_at IS NULL

AND rd.nip_baru = $1;
