begin;

alter table riwayat_penghargaan_umum
  drop column jenis_penghargaan,
  add column jenis_penghargaan_id int references ref_jenis_penghargaan(id);

commit;
