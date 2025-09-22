begin;

alter table riwayat_asesmen_nine_box RENAME TO riwayat_nine_box;

alter table riwayat_nine_box
  alter column tahun type varchar(4);

commit;
