begin;

comment on table ref_tingkat_pendidikan is 'Referensi referensi pendidikan';

alter table ref_jenis_kenaikan_pangkat rename to ref_jenis_kp;

alter table riwayat_asesmen rename to riwayat_assesmen;

alter table riwayat_hukuman_disiplin rename to riwayat_hukdis;

alter table riwayat_uji_kompetensi rename to riwayat_ujikom;

alter table orang_tua rename column tanggal_meninggal to tgl_meninggal;

alter table pegawai rename column tanggal_lahir to tgl_lahir;

alter table pegawai rename column tanggal_sk_cpns to tgl_sk_cpns;

alter table pegawai rename column tanggal_surat_dokter to tgl_surat_dokter;

alter table pegawai rename column tanggal_bebas_narkoba to tgl_bebas_narkoba;

alter table pegawai rename column tanggal_catatan_polisi to tgl_catatan_polisi;

alter table pegawai rename column tanggal_meninggal to tgl_meninggal;

alter table pegawai rename column tanggal_npwp to tgl_npwp;

commit;
