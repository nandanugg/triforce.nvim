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
