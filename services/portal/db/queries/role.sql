-- name: ListRoles :many
select
  r.id,
  r.nama,
  r.deskripsi,
  r.is_default,
  r.is_aktif,
  case
    when r.is_default then
      (
        select count(distinct nip) from "user" u
        where u.deleted_at is null
      )
    else
      (
        select count(1) from user_role ur
        where ur.role_id = r.id and ur.deleted_at is null
          and exists (select 1 from "user" u where u.nip = ur.nip and u.deleted_at is null)
      )
  end::int as jumlah_user
from role r
where r.deleted_at is null
order by r.nama
limit $1 offset $2;

-- name: CountRoles :one
select count(1)
from role
where deleted_at is null;

-- name: GetRole :one
select
  r.id,
  r.nama,
  r.deskripsi,
  r.is_default,
  r.is_aktif,
  case
    when r.is_default then
      (
        select count(distinct nip) from "user" u
        where u.deleted_at is null
      )
    else
      (
        select count(1) from user_role ur
        where ur.role_id = r.id and ur.deleted_at is null
          and exists (select 1 from "user" u where u.nip = ur.nip and u.deleted_at is null)
      )
  end::int as jumlah_user
from role r
where r.id = $1 and r.deleted_at is null;

-- name: CreateRole :one
insert into role (nama, deskripsi, is_default)
values ($1, $2, $3)
returning id;

-- name: UpdateRole :one
update role
set
  nama = coalesce(sqlc.narg('nama'), nama),
  deskripsi = coalesce(sqlc.narg('deskripsi'), deskripsi),
  is_default = coalesce(sqlc.narg('is_default'), is_default),
  is_aktif = coalesce(sqlc.narg('is_aktif'), is_aktif),
  updated_at = now()
where id = $1 and deleted_at is null
returning id;

-- name: ListRoleResourcePermissionsByRoleID :many
select
  id,
  resource_permission_id
from role_resource_permission
where role_id = $1 and deleted_at is null;

-- name: CreateRoleResourcePermissions :exec
insert into role_resource_permission (role_id, resource_permission_id)
select t.role_id, t.resource_permission_id
from (
  select
    @role_id::int2 as role_id,
    unnest(@resource_permission_ids::int4[]) as resource_permission_id
) as t
where not exists (
  select 1 from role_resource_permission as rrp
  where rrp.role_id = t.role_id
    and rrp.resource_permission_id = t.resource_permission_id
    and rrp.deleted_at is null
);

-- name: DeleteRoleResourcePermissions :exec
update role_resource_permission
set deleted_at = now()
where role_id = $1
  and resource_permission_id <> all(@exclude_resource_permission_ids::int4[])
  and deleted_at is null;

-- name: ListRolesByNIPs :many
select
  ur.nip,
  r.id,
  r.nama,
  r.is_default,
  r.is_aktif
from user_role ur
join role r on r.id = ur.role_id and r.is_default is false and r.deleted_at is null
where ur.nip = any(@nips::varchar[]) and ur.deleted_at is null
union all
select
  t.nip,
  r.id,
  r.nama,
  r.is_default,
  r.is_aktif
from (
  select unnest(@nips::varchar[]) as nip
) as t
join role r on r.is_default and r.deleted_at is null
order by nama;

-- name: CountRolesByIDs :one
select count(1) from role
where id = any(@ids::int2[]) and deleted_at is null;
