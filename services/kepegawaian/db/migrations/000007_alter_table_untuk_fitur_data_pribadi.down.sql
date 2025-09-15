begin;

alter table ref_kedudukan_hukum
  drop column is_pppk;

alter table pegawai
  drop constraint fk_pegawai_jabatan_instansi,
  drop constraint fk_pegawai_jabatan_instansi_real;

alter table pegawai
  alter column jabatan_instansi_id type int4 using jabatan_instansi_id::int4,
  alter column jabatan_instansi_real_id type int4 using jabatan_instansi_real_id::int4;

alter table ref_tingkat_pendidikan rename to tingkat_pendidikan;

alter table ref_pendidikan rename to pendidikan;

commit;
