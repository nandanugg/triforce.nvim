do $$
begin
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

alter table role drop column is_aktif;
