begin;

alter table "user"
  add column email varchar(255),
  add column nama varchar(255),
  add column unit_organisasi varchar(200);

commit;
