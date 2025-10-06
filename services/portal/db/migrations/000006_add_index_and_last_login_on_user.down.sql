begin;

drop index user_nip_idx;

drop index user_role_role_id_idx;

drop index role_is_default_idx;

alter table "user"
  drop column last_login_at;

commit;
