begin;

alter table ref_kedudukan_hukum
  add column is_pppk boolean default false;

alter table pegawai
  alter column jabatan_instansi_id type varchar(36),
  alter column jabatan_instansi_real_id type varchar(36);

alter table pegawai
  add constraint fk_pegawai_jabatan_instansi
    foreign key (jabatan_instansi_id) references ref_jabatan(kode_jabatan),
  add constraint fk_pegawai_jabatan_instansi_real
    foreign key (jabatan_instansi_real_id) references ref_jabatan(kode_jabatan);

alter table tingkat_pendidikan rename to ref_tingkat_pendidikan;

alter table pendidikan rename to ref_pendidikan;

commit;
