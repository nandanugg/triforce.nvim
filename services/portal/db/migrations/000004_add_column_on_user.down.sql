begin;

alter table "user"
  drop column email,
  drop column nama,
  drop column unit_organisasi;

commit;
