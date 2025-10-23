begin;

alter table pegawai rename column jml_pasangan to jml_istri;

alter table pegawai rename column satuan_kerja_id to satuan_kerja_kerja_id;

alter table ref_kelas_jabatan drop column deleted_at;

alter table riwayat_diklat rename column bkn_id to siasn_id;

alter table riwayat_kenaikan_gaji_berkala rename column tempat_lahir to birth_place;

alter table riwayat_kenaikan_gaji_berkala rename column tanggal_lahir to birth_date;

alter table riwayat_kenaikan_gaji_berkala rename column pendidikan_terakhir to last_education;

alter table riwayat_kenaikan_gaji_berkala rename column tanggal_lulus_pendidikan_terakhir to last_education_date;

alter table riwayat_kursus rename column bkn_id to siasn_id;

alter table ref_unit_kerja rename to unit_kerja;

alter table update_mandiri rename column tanggal_verifikasi to verifikasi_tgl;

alter table riwayat_surat_keputusan rename to file_digital_signature_riwayat;

alter table koreksi_surat_keputusan rename to file_digital_signature_corrector;

alter table surat_keputusan rename to file_digital_signature;

alter table log_request_surat_keputusan rename to log_digital_signature;

commit;
