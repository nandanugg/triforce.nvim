begin;

alter table riwayat_nine_box
  alter column tahun type int2 using tahun::int2;

alter table riwayat_nine_box RENAME TO riwayat_asesmen_nine_box;

commit;
