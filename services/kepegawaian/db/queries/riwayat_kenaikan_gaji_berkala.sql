-- name: ListRiwayatKenaikanGajiBerkala :many
SELECT
    rk.id,
    rg.id AS golongan_id,
    rg.nama AS golongan_nama,
    rg.nama_pangkat AS golongan_nama_pangkat,
    rk.no_sk,
    rk.tanggal_sk,
    rk.tmt_golongan,
    rk.masa_kerja_golongan_tahun,
    rk.masa_kerja_golongan_bulan,
    rk.tmt_sk AS tmt_kenaikan_gaji_berkala,
    rk.gaji_pokok,
    rk.jabatan,
    rk.tmt_jabatan,
    rk.pendidikan_terakhir AS pendidikan,
    rk.tanggal_lulus_pendidikan_terakhir AS tanggal_lulus,
    rk.kantor_pembayaran,
    rk.unit_kerja_induk_id,
    uk.nama_unor AS unit_kerja_induk,
    rk.pejabat
FROM
    riwayat_kenaikan_gaji_berkala rk
    JOIN pegawai p ON rk.pegawai_id = p.id
    LEFT JOIN ref_golongan rg ON rk.golongan_id = rg.id AND rg.deleted_at IS NULL
    LEFT JOIN ref_unit_kerja uk ON rk.unit_kerja_induk_id = uk.id AND uk.deleted_at IS NULL
WHERE
    p.nip_baru = $1
    AND p.deleted_at IS NULL
    AND rk.deleted_at IS NULL
ORDER BY
    rk.tmt_sk DESC
LIMIT
    $2
OFFSET
    $3;

-- name: CountRiwayatKenaikanGajiBerkala :one
SELECT
    COUNT(1)
FROM
    riwayat_kenaikan_gaji_berkala rk
    JOIN pegawai p ON rk.pegawai_id = p.id
WHERE
    p.nip_baru = $1
    AND p.deleted_at IS NULL
    AND rk.deleted_at IS NULL;

-- name: GetBerkasRiwayatKenaikanGajiBerkala :one
SELECT
    file_base64
FROM
    riwayat_kenaikan_gaji_berkala rk
    JOIN pegawai p ON rk.pegawai_id = p.id
WHERE
    p.nip_baru = $1
    AND rk.id = $2
    AND rk.deleted_at IS NULL;
