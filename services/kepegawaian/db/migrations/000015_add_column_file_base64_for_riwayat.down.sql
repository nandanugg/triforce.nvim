begin;

alter table riwayat_kursus
  drop column file_base64,
  drop column keterangan_berkas;

alter table riwayat_kgb
  drop column file_base64,
  drop column keterangan_berkas;

alter table riwayat_jabatan
  drop column file_base64,
  drop column keterangan_berkas;

commit;
