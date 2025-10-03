begin;

alter table role
  alter column id drop default,
  alter column id type int;

drop sequence role_id_seq;
create sequence role_id_seq as int owned by role.id;
select setval('role_id_seq', coalesce((select max(id) from role), 0) + 1, false);

update role set service = '' where service is null;

alter table role
  alter column id set default nextval('role_id_seq'),
  alter column service set not null,
  drop column deskripsi,
  drop column is_default;

comment on column role.service is null;

alter table user_role
  alter column role_id type int;

drop function permission_update_resource_permission_kode() cascade;

drop function resource_update_resource_permission_kode() cascade;

drop function resource_permission_set_kode() cascade;

drop table role_resource_permission;

drop table resource_permission;

drop table resource;

drop table permission;

commit;
