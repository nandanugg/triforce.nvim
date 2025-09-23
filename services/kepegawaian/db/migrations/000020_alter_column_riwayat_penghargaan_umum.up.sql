begin;

alter table riwayat_penghargaan_umum
  drop column jenis_penghargaan_id,
  add column jenis_penghargaan varchar(50);

commit;
