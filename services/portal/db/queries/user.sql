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
