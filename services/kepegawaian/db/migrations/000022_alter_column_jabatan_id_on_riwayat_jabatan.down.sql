begin;

alter table riwayat_jabatan
  drop constraint fk_riwayat_jabatan_jabatan_id;

alter table riwayat_jabatan
  alter column jabatan_id type int using jabatan_id::int;

commit;
