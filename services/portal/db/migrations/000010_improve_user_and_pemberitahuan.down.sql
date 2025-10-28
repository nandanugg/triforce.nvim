begin;

alter table "user"
  alter column email drop not null,
  alter column nama drop not null,
  add column unit_organisasi varchar(200);

alter table pemberitahuan
  alter column updated_at type date;

commit;
