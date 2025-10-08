begin;

alter table pindah_unit
  alter column nip type varchar(20);

alter table ref_jabatan
  drop column no,
  alter column nama_jabatan type varchar(200),
  alter column nama_jabatan_full type varchar(200),
  alter column nama_jabatan_bkn type varchar(200);

alter table ref_jenis_kenaikan_pangkat
  alter column nama type varchar(50);

alter table riwayat_asesmen
  alter column tahun_penilaian_id type smallint using tahun_penilaian_id::smallint;

alter table riwayat_diklat
  alter column institusi_penyelenggara type varchar(200),
  alter column no_sertifikat type varchar(100),
  alter column tahun_diklat type smallint,
  alter column durasi_jam type smallint,
  alter column nama_diklat type varchar(200);

alter table riwayat_diklat_fungsional
  alter column nama_kursus type varchar(200),
  alter column institusi_penyelenggara type varchar(200),
  alter column no_sertifikat type varchar(100),
  alter column instansi type varchar(200),
  alter column keterangan_berkas type varchar(200);

alter table riwayat_diklat_struktural
  alter column nomor type varchar(100);

alter table riwayat_golongan
  alter column id type integer using id::integer,
  alter column jenis_kp type varchar(50),
  alter column pangkat_nama type varchar(50),
  alter column sk_nomor type varchar(50),
  alter column jumlah_angka_kredit_utama type smallint,
  alter column jumlah_angka_kredit_tambahan type smallint;

create sequence riwayat_golongan_id_seq as integer owned by riwayat_golongan.id;
select setval('riwayat_golongan_id_seq', coalesce((select max(id) from riwayat_golongan), 0) + 1, false);

alter table riwayat_golongan
  alter column id set default nextval('riwayat_golongan_id_seq');

alter table riwayat_jabatan
  alter column catatan type varchar(200);

alter table riwayat_penghargaan_umum
  alter column deskripsi_penghargaan type varchar(100),
  alter column nama_penghargaan type varchar(200);

alter table riwayat_penugasan
  alter column deskripsi_jabatan type varchar(200),
  alter column nama_jabatan type varchar(200);

alter table riwayat_sertifikasi
  alter column nama_sertifikasi type varchar(255);

alter table unit_kerja
  alter column id type varchar(36),
  alter column kode_internal type varchar(36),
  alter column eselon_id type varchar(36),
  alter column cepat_kode type varchar(36),
  alter column diatasan_id type varchar(36),
  alter column instansi_id type varchar(36),
  alter column pemimpin_pns_id type varchar(36),
  alter column jenis_unor_id type varchar(36),
  alter column unor_induk type varchar(36),
  alter column eselon_1 type varchar(36),
  alter column eselon_2 type varchar(36),
  alter column eselon_3 type varchar(36),
  alter column eselon_4 type varchar(36),
  alter column jabatan_id type varchar(32);

commit;
