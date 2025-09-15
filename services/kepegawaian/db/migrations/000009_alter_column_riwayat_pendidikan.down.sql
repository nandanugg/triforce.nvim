begin;

alter table riwayat_pendidikan
  drop column file_base64,
  drop column keterangan_berkas;

commit;
