begin;

create index user_nip_idx on "user"(nip);

create index user_role_role_id_idx on user_role(role_id);

create index role_is_default_idx on role(is_default);

alter table "user"
  add column last_login_at timestamptz;

commit;
