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

-- name: CreateRiwayatKenaikanGajiBerkala :one
insert into riwayat_kenaikan_gaji_berkala (
    pegawai_id,
    tmt_sk,
    no_sk,
    pejabat,
    tanggal_sk,
    pegawai_nama,
    pegawai_nip,
    tempat_lahir,
    tanggal_lahir,
    n_gol_ruang,
    tmt_golongan,
    masa_kerja_golongan_tahun,
    masa_kerja_golongan_bulan,
    jabatan,
    tmt_jabatan,
    golongan_id,
    unit_kerja_induk_text,
    unit_kerja_induk_id,
    kantor_pembayaran,
    pendidikan_terakhir,
    tanggal_lulus_pendidikan_terakhir,
    created_at,
    updated_at,
    gaji_pokok
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18,
    $19,
    $20,
    $21,
    now(),
    now(),
    $22
) returning id;

-- name: UpdateRiwayatKenaikanGajiBerkala :execrows
update riwayat_kenaikan_gaji_berkala set
    tmt_sk = $1,
    no_sk = $2,
    pejabat = $3,
    tanggal_sk = $4,
    n_gol_ruang = $5,
    tmt_golongan = $6,
    masa_kerja_golongan_tahun = $7,
    masa_kerja_golongan_bulan = $8,
    jabatan = $9,
    tmt_jabatan = $10,
    golongan_id = $11,
    unit_kerja_induk_text = $12,
    unit_kerja_induk_id = $13,
    kantor_pembayaran = $14,
    pendidikan_terakhir = $15,
    tanggal_lulus_pendidikan_terakhir = $16,
    gaji_pokok = $17,
    updated_at = now()
where id = @id and pegawai_id = @pegawai_id and deleted_at is null;

-- name: DeleteRiwayatKenaikanGajiBerkala :execrows
update riwayat_kenaikan_gaji_berkala set
    deleted_at = now()
where id = @id and pegawai_id = @pegawai_id and deleted_at is null;

-- name: UploadBerkasRiwayatKenaikanGajiBerkala :execrows
update riwayat_kenaikan_gaji_berkala set
    file_base64 = @file_base64,
    updated_at = now()
where id = @id and pegawai_id = @pegawai_id and deleted_at is null;