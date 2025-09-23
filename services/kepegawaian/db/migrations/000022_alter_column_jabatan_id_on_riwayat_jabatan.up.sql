begin;

alter table riwayat_jabatan
  alter column jabatan_id type varchar(36);

alter table riwayat_jabatan
  add constraint fk_riwayat_jabatan_jabatan_id
    foreign key (jabatan_id) references ref_jabatan(kode_jabatan);

commit;
