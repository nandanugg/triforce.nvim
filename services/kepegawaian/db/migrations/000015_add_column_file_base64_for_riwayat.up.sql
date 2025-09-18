begin;

alter table riwayat_kursus
  add column file_base64 text,
  add column keterangan_berkas varchar(200);

alter table riwayat_kgb
  add column file_base64 text,
  add column keterangan_berkas varchar(200);

alter table riwayat_jabatan
  add column file_base64 text,
  add column keterangan_berkas varchar(200);

commit;
