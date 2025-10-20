-- name: GetUserNIPByIDAndSource :one
select nip from "user"
where id = $1 and source = $2 and deleted_at is null;

-- name: ListUserRoleByNIP :many
select distinct on (r.service)
  r.service,
  r.nama
from user_role ur
join role r on r.id = ur.role_id and r.deleted_at is null
where ur.nip = $1 and ur.deleted_at is null
order by r.service, ur.updated_at desc;

-- name: UpdateLastLoginAt :exec
update "user"
set last_login_at = now()
where id = $1 and source = $2;

-- name: ListUsersGroupByNIP :many
select
  u.nip,
  json_agg(
    json_build_object(
      'id', u.id,
      'source', u.source,
      'nama', u.nama,
      'email', u.email,
      'last_login_at', u.last_login_at
    )
    order by u.last_login_at desc nulls last
  ) as profiles
from "user" u
where u.deleted_at is null
  and (sqlc.narg('nip')::varchar is null or u.nip = sqlc.narg('nip')::varchar)
  and (
    sqlc.narg('role_id')::int2 is null
    or (
      select r.is_default from role r where r.id = sqlc.narg('role_id')::int2 and r.deleted_at is null
    ) is true
    or u.nip in (
      select ur.nip
      from user_role ur
      join role r on r.id = ur.role_id and r.deleted_at is null
      where ur.role_id = sqlc.narg('role_id')::int2 and ur.deleted_at is null
    )
  )
group by u.nip
limit $1 offset $2;

-- name: CountUsersGroupByNIP :one
select count(distinct u.nip)
from "user" u
where u.deleted_at is null
  and (sqlc.narg('nip')::varchar is null or u.nip = sqlc.narg('nip')::varchar)
  and (
    sqlc.narg('role_id')::int2 is null
    or (
      select r.is_default from role r where r.id = sqlc.narg('role_id')::int2 and r.deleted_at is null
    ) is true
    or u.nip in (
      select ur.nip
      from user_role ur
      join role r on r.id = ur.role_id and r.deleted_at is null
      where ur.role_id = sqlc.narg('role_id')::int2 and ur.deleted_at is null
    )
  );

-- name: GetUserGroupByNIP :one
select
  u.nip,
  json_agg(
    json_build_object(
      'id', u.id,
      'source', u.source,
      'nama', u.nama,
      'email', u.email,
      'last_login_at', u.last_login_at
    )
    order by u.last_login_at desc nulls last
  ) as profiles
from "user" u
where u.nip = $1 and u.deleted_at is null
group by u.nip;
