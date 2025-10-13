do $$
begin
  execute format($fmt$
    create or replace function resource_permission_set_kode()
    returns trigger
    language plpgsql
    set search_path = %I
    as $func$
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
    $func$;
  $fmt$, current_schema());

  execute format($fmt$
    create or replace function resource_update_resource_permission_kode()
    returns trigger
    language plpgsql
    set search_path = %I
    as $func$
      begin
        update resource_permission
        set resource_id = resource_id -- no-op update, just to regenerate kode.
        where resource_id = new.id;

        return null;
      end;
    $func$;
  $fmt$, current_schema());

  execute format($fmt$
    create or replace function permission_update_resource_permission_kode()
    returns trigger
    language plpgsql
    set search_path = %I
    as $func$
      begin
        update resource_permission
        set permission_id = permission_id -- no-op update, just to regenerate kode.
        where permission_id = new.id;

        return null;
      end;
    $func$;
  $fmt$, current_schema());

  execute format($fmt$
    create or replace function public.is_user_has_access(p_nip text, p_kode text)
    returns boolean
    language plpgsql
    security definer
    set search_path = %I
    as $func$
      begin
        return exists (
          select 1
          from role_resource_permission rrp
          join resource_permission rp on rp.id = rrp.resource_permission_id and rp.deleted_at is null
          join role r on r.id = rrp.role_id and r.deleted_at is null
          where rp.kode = p_kode
            and (r.is_default or rrp.role_id in (
              select role_id from user_role where nip = p_nip and deleted_at is null
            ))
            and rrp.deleted_at is null
        );
      end;
    $func$;
  $fmt$, current_schema());
end $$;
