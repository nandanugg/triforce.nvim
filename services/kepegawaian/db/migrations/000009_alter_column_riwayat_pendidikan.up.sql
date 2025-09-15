begin;

alter table riwayat_pendidikan
  alter column tahun_lulus type int2 using tahun_lulus::int2,
  add column file_base64 text,
  add column keterangan_berkas varchar(200);

commit;
