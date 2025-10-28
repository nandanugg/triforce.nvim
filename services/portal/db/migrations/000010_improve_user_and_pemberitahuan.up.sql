begin;

alter table "user"
  alter column email set not null,
  alter column nama set not null,
  drop column unit_organisasi;

alter table pemberitahuan
  alter column updated_at type timestamptz;

commit;
