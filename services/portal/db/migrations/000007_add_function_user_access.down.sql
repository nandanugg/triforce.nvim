begin;

create or replace function resource_permission_set_kode()
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

create or replace function resource_update_resource_permission_kode()
returns trigger as $$
  begin
    update resource_permission
    set resource_id = resource_id -- no-op update, just to regenerate kode.
    where resource_id = new.id;

    return null;
  end;
$$ language plpgsql;

create or replace function permission_update_resource_permission_kode()
returns trigger as $$
  begin
    update resource_permission
    set permission_id = permission_id -- no-op update, just to regenerate kode.
    where permission_id = new.id;

    return null;
  end;
$$ language plpgsql;

drop function public.is_user_has_access(text, text);

commit;
