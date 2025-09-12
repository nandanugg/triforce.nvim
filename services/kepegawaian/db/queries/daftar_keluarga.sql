-- name: ListOrangTuaByNip :many
SELECT
    ot.id,
    ot.hubungan,
    ot.nama,
    ot.tgl_meninggal,
    ot.no_dokumen AS nik,
    ot.agama_id,
    ot.akte_meninggal AS dokumen_pendukung,
    ra.nama AS agama_nama
FROM orang_tua ot
JOIN pegawai pg ON ot.pns_id = pg.pns_id AND pg.deleted_at is null
LEFT JOIN ref_agama ra ON ot.agama_id = ra.id AND ra.deleted_at is null
WHERE ot.deleted_at IS NULL
AND pg.nip_baru = $1;

-- name: ListPasanganByNip :many
SELECT
    p.id,
    p.pns,
    p.nama,
    p.tanggal_menikah,
    p.karsus AS nomor_karis,
    p.status,
    ra.nama AS agama_nama
FROM pasangan p
JOIN pegawai pg ON p.pns_id = pg.pns_id AND pg.deleted_at is null
LEFT JOIN ref_agama ra ON pg.agama_id = ra.id AND ra.deleted_at is null
WHERE p.deleted_at IS NULL
AND pg.nip_baru = $1;

-- name: ListAnakByNip :many
SELECT
    a.id,
    a.pasangan_id,
    a.nama,
    a.jenis_kelamin,
    a.tanggal_lahir,
    a.status_anak,
    ROW_NUMBER() OVER (
        PARTITION BY a.pns_id
        ORDER BY
            a.tanggal_lahir ASC NULLS LAST,
            a.id ASC  -- tie-breaker for same date or nulls
    ) AS anak_ke,
    -- a.status_sekolah,
    pas.nama AS nama_ibu_bapak
    -- '' AS dokumen_pendukung
FROM anak a
JOIN pegawai pg ON a.pns_id = pg.pns_id AND pg.deleted_at is null
LEFT JOIN pasangan pas ON a.pasangan_id = pas.id AND pas.deleted_at is null
WHERE a.deleted_at IS NULL
AND pg.nip_baru = $1
ORDER BY a.pns_id, anak_ke;
