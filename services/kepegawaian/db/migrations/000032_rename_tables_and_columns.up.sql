begin;

alter table pegawai rename column jml_istri to jml_pasangan;

alter table pegawai rename column satuan_kerja_kerja_id to satuan_kerja_id;

alter table ref_kelas_jabatan add column deleted_at timestamptz;

alter table riwayat_diklat rename column siasn_id to bkn_id;

alter table riwayat_kenaikan_gaji_berkala rename column birth_place to tempat_lahir;

alter table riwayat_kenaikan_gaji_berkala rename column birth_date to tanggal_lahir;

alter table riwayat_kenaikan_gaji_berkala rename column last_education to pendidikan_terakhir;

alter table riwayat_kenaikan_gaji_berkala rename column last_education_date to tanggal_lulus_pendidikan_terakhir;

alter table riwayat_kursus rename column siasn_id to bkn_id;

alter table unit_kerja rename to ref_unit_kerja;

alter table update_mandiri rename column verifikasi_tgl to tanggal_verifikasi;

alter table file_digital_signature_riwayat rename to riwayat_surat_keputusan;

alter table file_digital_signature_corrector rename to koreksi_surat_keputusan;

alter table file_digital_signature rename to surat_keputusan;

alter table log_digital_signature rename to log_request_surat_keputusan;

commit;
