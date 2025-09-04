-- name: GetParentsByEmployeeID :many
SELECT
    ot.nama,
    ot.tgl_meninggal,
    ot.no_dokumen AS nik,
    ra.nama AS agama,
    ot.hubungan,
    ot.pns_id
FROM
    orang_tua ot
LEFT JOIN
    ref_agama ra ON ot.agama_id = ra.id
WHERE
    ot.pns_id = $1;

-- name: GetSpouseByEmployeeID :one
SELECT
    p.nama,
    p.pns,
    p.karsus AS nomor_karis,
    p.status,
    p.pns_id
FROM
    pasangan p
WHERE
    p.pns_id = $1;

-- name: GetChildrenByEmployeeID :many
SELECT
    a.nama,
    a.jenis_kelamin,
    a.tanggal_lahir,
    pas.nama AS nama_pasangan,
    a.status_anak,
    a.pns_id
FROM
    anak a
LEFT JOIN
    pasangan pas ON a.pasangan_id = pas.id
WHERE
    a.pns_id = $1;

-- name: GetEmployeeFamilyData :one
SELECT
    p.nama AS nama_pegawai,
    p.id AS pns_id
FROM
    pegawai p
WHERE
    p.id = $1;
