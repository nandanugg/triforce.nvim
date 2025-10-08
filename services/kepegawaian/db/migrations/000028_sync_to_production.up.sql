begin;

alter table pindah_unit
  alter column nip type varchar(36);

alter table ref_jabatan
  add column no integer,
  alter column nama_jabatan type varchar(400),
  alter column nama_jabatan_full type varchar(400),
  alter column nama_jabatan_bkn type varchar(400);

update ref_jabatan set no = id;
create sequence ref_jabatan_no_seq as integer owned by ref_jabatan.no;
select setval('ref_jabatan_no_seq', coalesce((select max(no) from ref_jabatan), 0) + 1, false);

alter table ref_jabatan
  alter column no set not null,
  alter column no set default nextval('ref_jabatan_no_seq');

alter table ref_jenis_kenaikan_pangkat
  alter column nama type varchar(100);

alter table riwayat_asesmen
  alter column tahun_penilaian_id type varchar(10);

alter table riwayat_diklat
  alter column institusi_penyelenggara type varchar(600),
  alter column no_sertifikat type varchar(600),
  alter column tahun_diklat type integer,
  alter column durasi_jam type integer,
  alter column nama_diklat type varchar(700);

alter table riwayat_diklat_fungsional
  alter column nama_kursus type varchar(300),
  alter column institusi_penyelenggara type varchar(300),
  alter column no_sertifikat type varchar(200),
  alter column instansi type varchar(300),
  alter column keterangan_berkas type varchar(300);

alter table riwayat_diklat_struktural
  alter column nomor type varchar(300);

alter table riwayat_golongan
  alter column id drop default,
  alter column id type varchar(36),
  alter column jenis_kp type varchar(100),
  alter column pangkat_nama type varchar(100),
  alter column sk_nomor type varchar(100),
  alter column jumlah_angka_kredit_utama type integer,
  alter column jumlah_angka_kredit_tambahan type integer;

drop sequence riwayat_golongan_id_seq;

alter table riwayat_jabatan
  alter column catatan type varchar(250);

alter table riwayat_penghargaan_umum
  alter column deskripsi_penghargaan type varchar(1300),
  alter column nama_penghargaan type varchar(300);

alter table riwayat_penugasan
  alter column deskripsi_jabatan type varchar(3000),
  alter column nama_jabatan type varchar(400);

alter table riwayat_sertifikasi
  alter column nama_sertifikasi type varchar(300);

alter table unit_kerja
  alter column id type varchar(60),
  alter column kode_internal type varchar(60),
  alter column eselon_id type varchar(60),
  alter column cepat_kode type varchar(60),
  alter column diatasan_id type varchar(60),
  alter column instansi_id type varchar(60),
  alter column pemimpin_pns_id type varchar(60),
  alter column jenis_unor_id type varchar(60),
  alter column unor_induk type varchar(60),
  alter column eselon_1 type varchar(60),
  alter column eselon_2 type varchar(60),
  alter column eselon_3 type varchar(60),
  alter column eselon_4 type varchar(60),
  alter column jabatan_id type varchar(60);

commit;
