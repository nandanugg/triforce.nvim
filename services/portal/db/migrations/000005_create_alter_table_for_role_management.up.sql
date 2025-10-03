begin;

alter table role
  alter column id drop default,
  alter column id type smallint;

drop sequence role_id_seq;
create sequence role_id_seq as smallint owned by role.id;
select setval('role_id_seq', coalesce((select max(id) from role), 0) + 1, false);

alter table role
  alter column id set default nextval('role_id_seq'),
  alter column service drop not null,
  add column deskripsi varchar(255),
  add column is_default boolean not null default false;

comment on column role.service is 'deprecated';

alter table user_role
  alter column role_id type smallint;

create table resource (
  id smallserial primary key,
  service varchar(50) not null,
  kode varchar(50) not null,
  nama varchar(100) not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

create unique index resource_service_kode_unique_idx on resource(service, kode) where deleted_at is null;

create table permission (
  id smallserial primary key,
  kode varchar(50) not null,
  nama varchar(100) not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

create unique index permission_kode_unique_idx on permission(kode) where deleted_at is null;

create table resource_permission (
  id serial primary key,
  kode varchar(200),
  resource_id smallint not null references resource(id),
  permission_id smallint not null references permission(id),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

create unique index resource_permission_kode_unique_idx on resource_permission(kode) where deleted_at is null;
create index resource_permission_resource_id_idx on resource_permission(resource_id);
create index resource_permission_permission_id_idx on resource_permission(permission_id);

create table role_resource_permission (
  id serial primary key,
  role_id smallint not null references role(id),
  resource_permission_id int not null references resource_permission(id),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

create unique index role_resource_permission_unique_idx on role_resource_permission(role_id, resource_permission_id) where deleted_at is null;

-- trigger to set kode before insert/update on resource_permission.
create function resource_permission_set_kode()
returns trigger as $$
  declare
    kode varchar(200);
  begin
    select r.service || '.' || r.kode || '.' || p.kode into kode
    from resource r
    join permission p on p.id = new.permission_id and p.deleted_at is null
    where r.id = new.resource_id and r.deleted_at is null;

    new.kode = kode;
    return new;
  end;
$$ language plpgsql;

create trigger resource_permission_set_kode
before insert or update of resource_id, permission_id on resource_permission
for each row execute function resource_permission_set_kode();

-- trigger to update resource_permission.kode after update on resource
create function resource_update_resource_permission_kode()
returns trigger as $$
  begin
    update resource_permission
    set resource_id = resource_id -- no-op update, just to regenerate kode.
    where resource_id = new.id;

    return null;
  end;
$$ language plpgsql;

create trigger resource_update_resource_permission_kode
after update of service, kode, deleted_at on resource
for each row execute function resource_update_resource_permission_kode();

-- trigger to update resource_permission.kode after update on permission
create function permission_update_resource_permission_kode()
returns trigger as $$
  begin
    update resource_permission
    set permission_id = permission_id -- no-op update, just to regenerate kode.
    where permission_id = new.id;

    return null;
  end;
$$ language plpgsql;

create trigger permission_update_resource_permission_kode
after update of kode, deleted_at on permission
for each row execute function permission_update_resource_permission_kode();

commit;
