-- name: ListPendidikanFormal :many
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

