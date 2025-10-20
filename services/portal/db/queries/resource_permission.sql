-- name: ListResources :many
select
  id,
  nama
from resource
where deleted_at is null
order by kode
limit $1 offset $2;

-- name: CountResources :one
select count(1)
from resource
where deleted_at is null;

-- name: CountResourcePermissionsByIDs :one
select count(1)
from resource_permission
where kode is not null -- kode is not null is alias for resource.deleted_at is null and permission.deleted_at is null
  and id = any(@ids::int4[])
  and deleted_at is null;

-- name: ListResourcePermissionsByNip :many
select distinct rp.kode
from role_resource_permission rrp
join resource_permission rp on rp.id = rrp.resource_permission_id and rp.deleted_at is null
join role r on r.id = rrp.role_id and r.is_aktif is true and r.deleted_at is null
where rp.kode is not null -- rp.kode is not null is alias for resource.deleted_at is null and permission.deleted_at is null
  and (r.is_default or rrp.role_id in (
    select role_id from user_role where nip = $1 and deleted_at is null
  ))
  and rrp.deleted_at is null
order by rp.kode;

-- name: ListResourcePermissionsByResourceIDs :many
select
  rp.resource_id,
  rp.id,
  rp.kode,
  p.nama as nama_permission
from resource_permission rp
join permission p on p.id = rp.permission_id and p.deleted_at is null
where rp.resource_id = any(@resource_ids::int2[]) and rp.deleted_at is null
order by rp.kode;

-- name: ListResourcePermissionsByRoleID :many
select
  rp.id,
  rp.kode,
  r.nama as nama_resource,
  p.nama as nama_permission
from role_resource_permission rrp
join resource_permission rp on rp.id = rrp.resource_permission_id and rp.deleted_at is null
join resource r on r.id = rp.resource_id and r.deleted_at is null
join permission p on p.id = rp.permission_id and p.deleted_at is null
where rrp.role_id = $1 and rrp.deleted_at is null
order by rp.kode;
