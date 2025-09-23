begin;

comment on table ref_tingkat_pendidikan is 'Referensi tingkat pendidikan';

alter table ref_jenis_kp rename to ref_jenis_kenaikan_pangkat;

alter table riwayat_assesmen rename to riwayat_asesmen;

alter table riwayat_hukdis rename to riwayat_hukuman_disiplin;

alter table riwayat_ujikom rename to riwayat_uji_kompetensi;

alter table orang_tua rename column tgl_meninggal to tanggal_meninggal;

alter table pegawai rename column tgl_lahir to tanggal_lahir;

alter table pegawai rename column tgl_sk_cpns to tanggal_sk_cpns;

alter table pegawai rename column tgl_surat_dokter to tanggal_surat_dokter;

alter table pegawai rename column tgl_bebas_narkoba to tanggal_bebas_narkoba;

alter table pegawai rename column tgl_catatan_polisi to tanggal_catatan_polisi;

alter table pegawai rename column tgl_meninggal to tanggal_meninggal;

alter table pegawai rename column tgl_npwp to tanggal_npwp;

commit;
