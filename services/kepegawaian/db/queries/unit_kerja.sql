-- name: ListUnitKerjaByNamaOrInduk :many
SELECT id, nama_unor
from unit_kerja
WHERE
    (CASE WHEN @nama::varchar = '' THEN true ELSE nama_unor ilike @nama::varchar || '%' END)
    AND (CASE WHEN @unor_induk::varchar = '' THEN true ELSE unor_induk = @unor_induk::varchar END)
    AND deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountUnitKerja :one
SELECT COUNT(1) FROM unit_kerja
WHERE
    (CASE WHEN @nama::varchar = '' THEN true ELSE nama_unor ilike @nama::varchar || '%' END)
    AND (CASE WHEN @unor_induk::varchar = '' THEN true ELSE unor_induk = @unor_induk::varchar END)
    AND deleted_at IS NULL;

-- name: ListUnitKerjaHierarchy :many
with recursive unit_kerja_path as (
    -- anchor
    select uk.id, uk.nama_unor, uk.diatasan_id, 1 as depth
    from unit_kerja uk
    where uk.id = $1 and uk.deleted_at is null

    union all

    -- recursive
    select uk.id, uk.nama_unor, uk.diatasan_id, ukp.depth + 1
    from unit_kerja uk
    join unit_kerja_path ukp on uk.id = ukp.diatasan_id
    where ukp.depth < 10 and uk.deleted_at is null
)
select id, nama_unor from unit_kerja_path;

-- name: ListUnitKerjaHierarchyByNIP :many
with recursive unit_kerja_path as (
    -- anchor
    select uk.id, uk.nama_unor, uk.diatasan_id, 1 as depth
    from unit_kerja uk
    where uk.id = (SELECT unor_id FROM pegawai where nip_baru = @nip::varchar LIMIT 1) and uk.deleted_at is null

    union all

    -- recursive
    select uk.id, uk.nama_unor, uk.diatasan_id, ukp.depth + 1
    from unit_kerja uk
    join unit_kerja_path ukp on uk.id = ukp.diatasan_id
    where ukp.depth < 10 and uk.deleted_at is null
)
select id, nama_unor from unit_kerja_path;

-- name: ListUnitKerjaLengkapByIDs :many
with recursive unit_kerja_path as (
    -- anchor
    select 
        uk.id,
        uk.id as start_id,
        uk.diatasan_id,
        uk.nama_unor::text as path,
        1 as depth
    from unit_kerja uk
    where uk.id = ANY(sqlc.arg(ids)::varchar[]) AND uk.deleted_at is null

    union all

    -- recursive
    select 
        uk.id,
        ukp.start_id,
        uk.diatasan_id,
        ukp.path || ' - ' || uk.nama_unor, 
        ukp.depth + 1
    from unit_kerja uk
    join unit_kerja_path ukp
      on uk.id = ukp.diatasan_id
    where ukp.depth < 10 and uk.deleted_at is null
)

select 
    start_id as id,
    path as nama_unor_lengkap
from (
    select *, row_number() over (partition by start_id order by depth desc) as rn
    from unit_kerja_path
) t
where rn = 1;

-- name: ListAkarUnitKerja :many
SELECT 
    uk.id, 
    uk.nama_unor,
    EXISTS (
        SELECT 1 
        FROM unit_kerja uk2
        WHERE 
            uk2.diatasan_id = uk.id
            AND uk2.deleted_at IS NULL
    ) as has_anak
FROM unit_kerja uk
WHERE
    uk.diatasan_id IS NULL
    AND uk.deleted_at IS NULL
ORDER BY "order"
LIMIT $1 OFFSET $2;

-- name: CountAkarUnitKerja :one
SELECT COUNT(1) FROM unit_kerja
WHERE
    diatasan_id IS NULL
    AND deleted_at IS NULL;

-- name: ListUnitKerjaByDiatasanID :many
SELECT 
    uk.id,
    uk.nama_unor,
    EXISTS (
        SELECT 1 
        FROM unit_kerja uk2
        WHERE 
            uk2.diatasan_id = uk.id
            AND uk2.deleted_at IS NULL
    ) as has_anak
FROM unit_kerja uk
WHERE
    uk.diatasan_id = sqlc.arg(diatasan_id)
    AND uk.deleted_at IS NULL
ORDER BY "order"
LIMIT $1 OFFSET $2;

-- name: CountUnitKerjaByDiatasanID :one
SELECT COUNT(1) FROM unit_kerja
WHERE
    diatasan_id = sqlc.arg(diatasan_id)
    AND deleted_at IS NULL;

-- name: GetUnitKerja :one
SELECT
    uk.id,
    uk."no",
    uk.kode_internal,
    uk.nama_unor as nama,
    uk.eselon_id,
    uk.cepat_kode,
    uk.nama_jabatan,
    uk.nama_pejabat,
    uk.diatasan_id,
    uk.instansi_id,
    uk.pemimpin_pns_id,
    uk.jenis_unor_id,
    uk.unor_induk,
    uk.jumlah_ideal_staff,
    uk."order",
    uk.is_satker,
    uk.eselon_1,
    uk.eselon_2,
    uk.eselon_3,
    uk.eselon_4,
    uk.expired_date,
    uk.keterangan,
    uk.jenis_satker,
    uk.abbreviation,
    uk.unor_induk_penyetaraan,
    uk.jabatan_id,
    uk.waktu,
    uk.peraturan,
    uk.remark,
    uk.aktif,
    uk.eselon_nama,
    ukd.nama_unor as nama_diatasan,
    ukui.nama_unor as nama_unor_induk
FROM unit_kerja uk
LEFT JOIN unit_kerja ukd ON uk.diatasan_id = ukd.id AND ukd.deleted_at IS NULL
LEFT JOIN unit_kerja ukui ON uk.unor_induk = ukui.id AND ukui.deleted_at IS NULL
WHERE uk.id = @id::varchar AND uk.deleted_at IS NULL;

-- name: CreateUnitKerja :one
INSERT INTO unit_kerja (
  diatasan_id,
  id,
  nama_unor,
  kode_internal,
  nama_jabatan,
  pemimpin_pns_id,
  nama_pejabat,
  is_satker,
  unor_induk,
  expired_date,
  keterangan,
  abbreviation,
  waktu,
  jenis_satker,
  peraturan
) VALUES (
  @diatasan_id,
  @id,
  @nama,
  @kode_internal,
  @nama_jabatan,
  @pemimpin_pns_id,
  @nama_pejabat,
  @is_satker,
  @unor_induk,
  @expired_date,
  @keterangan,
  @abbreviation,
  @waktu,
  @jenis_satker,
  @peraturan
)
RETURNING
  id,
  "no",
  kode_internal,
  nama_unor as nama,
  eselon_id,
  cepat_kode,
  nama_jabatan,
  nama_pejabat,
  diatasan_id,
  instansi_id,
  pemimpin_pns_id,
  jenis_unor_id,
  unor_induk,
  jumlah_ideal_staff,
  "order",
  is_satker,
  eselon_1,
  eselon_2,
  eselon_3,
  eselon_4,
  expired_date,
  keterangan,
  jenis_satker,
  abbreviation,
  unor_induk_penyetaraan,
  jabatan_id,
  waktu,
  peraturan,
  remark,
  aktif,
  eselon_nama;

-- name: UpdateUnitKerja :one
UPDATE unit_kerja
SET
    diatasan_id = @diatasan_id,
    nama_unor = @nama,
    kode_internal = @kode_internal,
    nama_jabatan = @nama_jabatan,
    pemimpin_pns_id = @pemimpin_pns_id,
    nama_pejabat = @nama_pejabat,
    is_satker = @is_satker,
    unor_induk = @unor_induk,
    expired_date = @expired_date,
    keterangan = @keterangan,
    abbreviation = @abbreviation,
    waktu = @waktu,
    jenis_satker = @jenis_satker,
    peraturan = @peraturan,
    updated_at = now()
WHERE id = @id AND deleted_at IS NULL
RETURNING
    id,
    "no",
    kode_internal,
    nama_unor as nama,
    eselon_id,
    cepat_kode,
    nama_jabatan,
    nama_pejabat,
    diatasan_id,
    instansi_id,
    pemimpin_pns_id,
    jenis_unor_id,
    unor_induk,
    jumlah_ideal_staff,
    "order",
    is_satker,
    eselon_1,
    eselon_2,
    eselon_3,
    eselon_4,
    expired_date,
    keterangan,
    jenis_satker,
    abbreviation,
    unor_induk_penyetaraan,
    jabatan_id,
    waktu,
    peraturan,
    remark,
    aktif,
    eselon_nama;

-- name: DeleteUnitKerja :execrows
UPDATE unit_kerja
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
