-- name: ListPendidikanFormal :many
SELECT 
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
FROM riwayat_pendidikan rp
    LEFT JOIN tingkat_pendidikan tk ON tk.id = rp.tingkat_pendidikan_id AND rp.deleted_at IS NULL
    LEFT JOIN pendidikan pend ON rp.pendidikan_id = pend.id AND pend.deleted_at IS NULL
JOIN pegawai p ON p.pns_id = rp.pns_id_3
WHERE p.nip_baru = $1
ORDER BY rp.tahun_lulus ASC
;
