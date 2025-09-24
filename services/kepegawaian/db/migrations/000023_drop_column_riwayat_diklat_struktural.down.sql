begin;

alter table riwayat_diklat_struktural
  add column institusi_penyelenggara varchar(200);

commit;
