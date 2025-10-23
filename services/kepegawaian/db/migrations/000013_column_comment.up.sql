-- Add table comments

COMMENT ON TABLE orang_tua IS 'Orang tua pegawai';
COMMENT ON TABLE anak IS 'Anak pegawai';
COMMENT ON TABLE pasangan IS 'Pasangan pegawai';
COMMENT ON TABLE pegawai IS 'Data utama pegawai';
COMMENT ON TABLE pindah_unit IS 'Riwayat perpindahan unit kerja';
COMMENT ON TABLE ref_agama IS 'Referensi agama';
COMMENT ON TABLE ref_golongan IS 'Referensi golongan';
COMMENT ON TABLE ref_instansi IS 'Referensi instansi';
COMMENT ON TABLE ref_jabatan IS 'Referensi jabatan';
COMMENT ON TABLE ref_jenis_diklat IS 'Referensi jenis diklat';
COMMENT ON TABLE ref_jenis_diklat_fungsional IS 'Referensi jenis diklat fungsional';
COMMENT ON TABLE ref_jenis_diklat_struktural IS 'Referensi jenis diklat struktural';
COMMENT ON TABLE ref_jenis_hukuman IS 'Referensi jenis hukuman disiplin';
COMMENT ON TABLE ref_jenis_jabatan IS 'Referensi jenis jabatan';
COMMENT ON TABLE ref_jenis_kawin IS 'Referensi status pernikahan';
COMMENT ON TABLE ref_jenis_kp IS 'Referensi jenis kenaikan pangkat';
COMMENT ON TABLE ref_jenis_pegawai IS 'Referensi jenis pegawai';
COMMENT ON TABLE ref_jenis_penghargaan IS 'Referensi jenis penghargaan';
COMMENT ON TABLE ref_jenis_satker IS 'Referensi jenis satuan kerja';
COMMENT ON TABLE ref_kedudukan_hukum IS 'Referensi kedudukan hukum';
COMMENT ON TABLE ref_kelas_jabatan IS 'Referensi kelas jabatan';
COMMENT ON TABLE ref_kpkn IS 'Referensi kantor pelayanan kekayaan negara';
COMMENT ON TABLE ref_lokasi IS 'Referensi lokasi';
COMMENT ON TABLE ref_pendidikan IS 'Referensi pendidikan';
COMMENT ON TABLE ref_tingkat_pendidikan IS 'Referensi referensi pendidikan';
COMMENT ON TABLE riwayat_assesmen IS 'Riwayat asesmen pegawai';
COMMENT ON TABLE riwayat_diklat IS 'Riwayat diklat pegawai';
COMMENT ON TABLE riwayat_diklat_fungsional IS 'Riwayat diklat fungsional pegawai';
COMMENT ON TABLE riwayat_diklat_struktural IS 'Riwayat diklat struktural pegawai';
COMMENT ON TABLE riwayat_golongan IS 'Riwayat golongan pegawai';
COMMENT ON TABLE riwayat_hukdis IS 'Riwayat hukuman disiplin pegawai';
COMMENT ON TABLE riwayat_jabatan IS 'Riwayat jabatan pegawai';
COMMENT ON TABLE riwayat_kgb IS 'Riwayat kenaikan gaji berkala pegawai';
COMMENT ON TABLE riwayat_kinerja IS 'Riwayat kinerja pegawai';
COMMENT ON TABLE riwayat_kursus IS 'Riwayat kursus pegawai';
COMMENT ON TABLE riwayat_nine_box IS 'Riwayat asesmen nine box pegawai';
COMMENT ON TABLE riwayat_pendidikan IS 'Riwayat pendidikan pegawai';
COMMENT ON TABLE riwayat_penghargaan_umum IS 'Riwayat penghargaan umum pegawai';
COMMENT ON TABLE riwayat_penugasan IS 'Riwayat penugasan pegawai';
COMMENT ON TABLE riwayat_pindah_unit_kerja IS 'Riwayat pindah unit kerja pegawai';
COMMENT ON TABLE riwayat_sertifikasi IS 'Riwayat sertifikasi pegawai';
COMMENT ON TABLE riwayat_ujikom IS 'Riwayat uji kompetensi pegawai';
COMMENT ON TABLE unit_kerja IS 'Referensi referensi unit kerja';
COMMENT ON TABLE update_mandiri IS 'Riwayat pembaruan data secara mandiri oleh pegawai';

-- Add column comments

-- Comments for table: anak
COMMENT ON COLUMN anak.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN anak.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN anak.id IS 'id data anak';
COMMENT ON COLUMN anak.jenis_kelamin IS 'Kode jenis kelamin anak, M: laki-laki, F: perempuan';
COMMENT ON COLUMN anak.nama IS 'Nama anak';
COMMENT ON COLUMN anak.nip IS 'NIP pegawai';
COMMENT ON COLUMN anak.pasangan_id IS 'Kaitan ke data pasangan (jika anak dikaitkan ke pasangan tertentu)';
COMMENT ON COLUMN anak.pns_id IS 'Referensi ke pegawai.pns_id';
COMMENT ON COLUMN anak.status_anak IS 'Status anak, 1: kandung, 2: angkat';
COMMENT ON COLUMN anak.tanggal_lahir IS 'Tanggal lahir anak';
COMMENT ON COLUMN anak.tempat_lahir IS 'Tempat lahir anak';
COMMENT ON COLUMN anak.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: orang_tua
COMMENT ON COLUMN orang_tua.agama_id IS 'id agama orang tua (rujuk ref_agama)';
COMMENT ON COLUMN orang_tua.akte_meninggal IS 'Nomor akte meninggal orang tua';
COMMENT ON COLUMN orang_tua.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN orang_tua.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN orang_tua.email IS 'Alamat email orang tua (bila ada)';
COMMENT ON COLUMN orang_tua.gelar_belakang IS 'Gelar di belakang nama orang tua';
COMMENT ON COLUMN orang_tua.gelar_depan IS 'Gelar di depan nama orang tua';
COMMENT ON COLUMN orang_tua.hubungan IS 'Kode hubungan, 1: ayah, 2: ibu';
COMMENT ON COLUMN orang_tua.id IS 'id data orang tua';
COMMENT ON COLUMN orang_tua.jenis_dokumen IS 'Jenis dokumen identitas, enum: KTP, PASPOR';
COMMENT ON COLUMN orang_tua.nama IS 'Nama lengkap orang tua';
COMMENT ON COLUMN orang_tua.nip IS 'NIP pegawai';
COMMENT ON COLUMN orang_tua.no_dokumen IS 'Nomor dokumen identitas';
COMMENT ON COLUMN orang_tua.pns_id IS 'Referensi ke pegawai.pns_id';
COMMENT ON COLUMN orang_tua.tanggal_lahir IS 'Tanggal lahir orang tua';
COMMENT ON COLUMN orang_tua.tempat_lahir IS 'Tempat lahir orang tua';
COMMENT ON COLUMN orang_tua.tgl_meninggal IS 'Tanggal meninggal orang tua';
COMMENT ON COLUMN orang_tua.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: pasangan
COMMENT ON COLUMN pasangan.akte_cerai IS 'Nomor/berkas akta cerai';
COMMENT ON COLUMN pasangan.akte_meninggal IS 'Nomor/berkas akta kematian';
COMMENT ON COLUMN pasangan.akte_nikah IS 'Nomor/berkas akta nikah';
COMMENT ON COLUMN pasangan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN pasangan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN pasangan.hubungan IS 'Kode hubungan, 1: istri, 2: suami';
COMMENT ON COLUMN pasangan.id IS 'id data pasangan';
COMMENT ON COLUMN pasangan.karsus IS 'Nomor kartu suami istri';
COMMENT ON COLUMN pasangan.nama IS 'Nama lengkap pasangan';
COMMENT ON COLUMN pasangan.nip IS 'NIP pegawai';
COMMENT ON COLUMN pasangan.pns IS 'Penanda apakah pasangan juga PNS';
COMMENT ON COLUMN pasangan.pns_id IS 'Referensi ke pegawai.pns_id';
COMMENT ON COLUMN pasangan.status IS 'Status hubungan saat ini, 1: menikah, 2: cerai, 3: jada/duda';
COMMENT ON COLUMN pasangan.tanggal_cerai IS 'Tanggal perceraian (jika ada)';
COMMENT ON COLUMN pasangan.tanggal_lahir IS 'Tanggal lahir pasangan';
COMMENT ON COLUMN pasangan.tanggal_menikah IS 'Tanggal pernikahan';
COMMENT ON COLUMN pasangan.tanggal_meninggal IS 'Tanggal wafat pasangan (jika ada)';
COMMENT ON COLUMN pasangan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: pegawai
COMMENT ON COLUMN pegawai.agama_id IS 'Kode agama (rujuk ref_agama)';
COMMENT ON COLUMN pegawai.akte_kelahiran IS 'Nomor akte kelahiran';
COMMENT ON COLUMN pegawai.akte_meninggal IS 'Nomor akte meninggal';
COMMENT ON COLUMN pegawai.alamat IS 'Alamat domisili pegawai';
COMMENT ON COLUMN pegawai.bpjs IS 'Nomor kepesertaan BPJS';
COMMENT ON COLUMN pegawai.bup IS 'Batas usia pensiun';
COMMENT ON COLUMN pegawai.created_at IS 'Waktu perekaman data dibuat';
COMMENT ON COLUMN pegawai.created_by IS 'id user yang memasukkan data pegawai';
COMMENT ON COLUMN pegawai.deleted_at IS 'Waktu penghapusan lunak (soft delete) bila ada';
COMMENT ON COLUMN pegawai.email IS 'Alamat surat elektronik pribadi';
COMMENT ON COLUMN pegawai.email_dikbud IS 'Alamat surat elektronik untuk kepentingan pekerjaan';
COMMENT ON COLUMN pegawai.email_dikbud_bak IS 'Alamat surat elektronik backup untuk kepentingan pekerjaan'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN pegawai.foto IS 'Lokasi/URL berkas foto pegawai';
COMMENT ON COLUMN pegawai.gelar_belakang IS 'Gelar akademik/jabatan di belakang nama';
COMMENT ON COLUMN pegawai.gelar_depan IS 'Gelar akademik/jabatan di depan nama';
COMMENT ON COLUMN pegawai.gol_awal_id IS 'Golongan awal saat pengangkatan (rujuk ref_golongan)';
COMMENT ON COLUMN pegawai.gol_id IS 'Golongan terakhir/aktif (rujuk ref_golongan)';
COMMENT ON COLUMN pegawai.golongan_darah IS 'Golongan darah';
COMMENT ON COLUMN pegawai.id IS 'Identitas numerik unik baris data pegawai';
COMMENT ON COLUMN pegawai.instansi_induk_id IS 'id instansi induk pegawai (rujuk ref_instansi)';
COMMENT ON COLUMN pegawai.instansi_induk_nama IS 'Nama instansi induk pegawai';
COMMENT ON COLUMN pegawai.instansi_kerja_id IS 'id instansi tempat bekerja (rujuk ref_instansi)';
COMMENT ON COLUMN pegawai.instansi_kerja_nama IS 'Nama instansi tempat bekerja pegawai';
-- TODO: hapus jika sudah tidak digunakan - is_dosen
COMMENT ON COLUMN pegawai.jabatan_id IS 'id jabatan pegawai (rujuk ref_jabatan)';
COMMENT ON COLUMN pegawai.jabatan_instansi_id IS 'id jabatan instansi pegawai (rujuk ref_jabatan)';
COMMENT ON COLUMN pegawai.jabatan_instansi_nama IS 'Nama jabatan instansi pegawai';
COMMENT ON COLUMN pegawai.jabatan_instansi_real_id IS 'id jabatan instansi pegawai (rujuk ref_jabatan)';
COMMENT ON COLUMN pegawai.jabatan_nama IS 'Nama jabatan pegawai';
COMMENT ON COLUMN pegawai.jabatan_ppnpn IS 'Nama jabatan Pegawai Pemerintah Non Pegawai Negeri';
COMMENT ON COLUMN pegawai.jenis_jabatan_id IS 'id jenis jabatan';
COMMENT ON COLUMN pegawai.jenis_jabatan_nama IS 'Nama jenis jabatan';
COMMENT ON COLUMN pegawai.jenis_kawin_id IS 'Status perkawinan (rujuk ref_jenis_kawin)';
COMMENT ON COLUMN pegawai.jenis_kelamin IS 'Kode jenis kelamin, M: laki-laki, F: perempuan';
COMMENT ON COLUMN pegawai.jenis_pegawai_id IS 'id jenis pegawai (PNS/PPPK/dll; rujuk ref_jenis_pegawai)';
COMMENT ON COLUMN pegawai.jml_anak IS 'Jumlah anak yang tercatat';
COMMENT ON COLUMN pegawai.jml_istri IS 'Jumlah pasangan';
COMMENT ON COLUMN pegawai.kartu_asn IS 'Nomor kartu ASN';
COMMENT ON COLUMN pegawai.kartu_pegawai IS 'Nomor kartu pegawai';
COMMENT ON COLUMN pegawai.kedudukan_hukum_id IS 'id kedudukan hukum (rujuk ref_kedudukan_hukum)';
COMMENT ON COLUMN pegawai.ket IS 'Keterangan tambahan terhadap pegawai';
COMMENT ON COLUMN pegawai.kk IS 'Nomor kartu keluarga';
 -- TODO: hapus jika sudah tidak digunakan - kodecepat
COMMENT ON COLUMN pegawai.kpkn_id IS 'id KPPN/KPKN pembayaran gaji (rujuk ref_kpkn)';
COMMENT ON COLUMN pegawai.kpkn_nama IS 'Nama KPPN/KPKN pembayaran gaji';
COMMENT ON COLUMN pegawai.lokasi_kerja IS 'Nama lokasi kerja';
COMMENT ON COLUMN pegawai.lokasi_kerja_id IS 'id lokasi kerja (rujuk ref_lokasi)';
COMMENT ON COLUMN pegawai.mk_bulan IS 'Masa kerja bulan';
COMMENT ON COLUMN pegawai.mk_bulan_swasta IS 'Masa kerja bulan di swasta, sebelum menjadi ASN';
COMMENT ON COLUMN pegawai.mk_tahun IS 'Masa kerja tahun';
COMMENT ON COLUMN pegawai.mk_tahun_swasta IS 'Masa kerja tahun di swasta, sebelum menjadi ASN';
COMMENT ON COLUMN pegawai.nama IS 'Nama lengkap pegawai';
 -- TODO: hapus jika sudah tidak digunakan - nidn
COMMENT ON COLUMN pegawai.nik IS 'Nomor Induk Kependudukan';
COMMENT ON COLUMN pegawai.nip_baru IS 'Nomor Induk Pegawai format baru (20 digit)';
COMMENT ON COLUMN pegawai.nip_lama IS 'Nomor Induk Pegawai format lama';
COMMENT ON COLUMN pegawai.no_askes IS 'Nomor ASKES (jika tersedia/legacy)';
COMMENT ON COLUMN pegawai.no_bebas_narkoba IS 'Nomor Surat Keterangan Bebas Narkoba';
COMMENT ON COLUMN pegawai.no_catatan_polisi IS 'Nomor Surat Catatan Kelakukan Baik dari kepolisian';
COMMENT ON COLUMN pegawai.no_darurat IS 'Nomor telepon yang dapat dihubungi dalam keadaan darurat';
COMMENT ON COLUMN pegawai.no_hp IS 'Nomor telepon seluler pegawai';
COMMENT ON COLUMN pegawai.no_sk_cpns IS 'Nomor SK pengangkatan CPNS';
COMMENT ON COLUMN pegawai.no_sk_pemberhentian IS 'Nomor SK pemberhentian dari PNS';
COMMENT ON COLUMN pegawai.no_surat_dokter IS 'Nomor surat pemeriksaan kesehatan';
COMMENT ON COLUMN pegawai.no_taspen IS 'Nomor Taspen';
COMMENT ON COLUMN pegawai.npwp IS 'Nomor Pokok Wajib Pajak';
COMMENT ON COLUMN pegawai.pendidikan_id IS 'id pendidikan (rujuk ref_pendidikan)';
COMMENT ON COLUMN pegawai.pns_id IS 'id pegawai negeri sipil (UUID) yang menjadi kunci rujukan antar tabel';
COMMENT ON COLUMN pegawai.satuan_kerja_induk_id IS 'id satuan kerja induk pegawai';
COMMENT ON COLUMN pegawai.satuan_kerja_induk_nama IS 'Nama satuan kerja induk pegawai';
COMMENT ON COLUMN pegawai.satuan_kerja_kerja_id IS 'id satuan kerja pegawai';
COMMENT ON COLUMN pegawai.satuan_kerja_nama IS 'Nama satuan kerja pegawai';
COMMENT ON COLUMN pegawai.status_cpns_pns IS 'Status CPNS/PNS';
COMMENT ON COLUMN pegawai.status_hidup IS 'Status hidup pegawai';
COMMENT ON COLUMN pegawai.status_pegawai IS 'Status pegawai, 1: pns, 2: honorer';
COMMENT ON COLUMN pegawai.status_pegawai_backup IS 'Status pegawai backup'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN pegawai.tahun_lulus IS 'Tahun kelulusan pendidikan terakhir';
COMMENT ON COLUMN pegawai.tempat_lahir IS 'Nama tempat lahir berdasarkan referensi ref_lokasi';
COMMENT ON COLUMN pegawai.tempat_lahir_id IS 'id tempat lahir (rujuk ref_lokasi)';
COMMENT ON COLUMN pegawai.tempat_lahir_nama IS 'Nama tempat lahir (teks bebas)';
COMMENT ON COLUMN pegawai.tgl_bebas_narkoba IS 'Tanggal Surat Keterangan Bebas Narkoba';
COMMENT ON COLUMN pegawai.tgl_catatan_polisi IS 'Tanggal Surat Catatan Kelakukan Baik dari kepolisian';
COMMENT ON COLUMN pegawai.tgl_lahir IS 'Tanggal lahir pegawai';
COMMENT ON COLUMN pegawai.tgl_meninggal IS 'Tanggal meninggal pegawai';
COMMENT ON COLUMN pegawai.tgl_npwp IS 'Tanggal terbit NPWP';
COMMENT ON COLUMN pegawai.tgl_sk_cpns IS 'Tanggal SK pengangkatan CPNS';
COMMENT ON COLUMN pegawai.tgl_surat_dokter IS 'Tanggal surat pemeriksaan kesehatan';
COMMENT ON COLUMN pegawai.tingkat_pendidikan_id IS 'Tingkat pendidikan terakhir (rujuk tingkat_pendidikan)';
COMMENT ON COLUMN pegawai.tmt_cpns IS 'Tanggal mulai tugas (CPNS)';
COMMENT ON COLUMN pegawai.tmt_golongan IS 'Tanggal mulai berlaku golongan saat ini';
COMMENT ON COLUMN pegawai.tmt_jabatan IS 'Tanggal mulai jabatan';
COMMENT ON COLUMN pegawai.tmt_pensiun IS 'Tanggal perkiraan/penetapan pensiun (BUP)';
COMMENT ON COLUMN pegawai.tmt_pns IS 'Tanggal mulai tugas (PNS)';
COMMENT ON COLUMN pegawai.unor_id IS 'Unit organisasi/kerja (rujuk unit_kerja)';
COMMENT ON COLUMN pegawai.unor_induk_id IS 'Unit organisasi/kerja induk (rujuk unit_kerja)';
COMMENT ON COLUMN pegawai.updated_at IS 'Waktu terakhir data diperbarui';
COMMENT ON COLUMN pegawai.updated_by IS 'id user yang memperbarui data pegawai';

-- Comments for table: pindah_unit
COMMENT ON COLUMN pindah_unit.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN pindah_unit.created_by IS 'id user yang membuat data pindah unit';
COMMENT ON COLUMN pindah_unit.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN pindah_unit.file_sk IS 'Lokasi penyimpanan berkas pindah unit';
COMMENT ON COLUMN pindah_unit.id IS 'id pindah unit';
COMMENT ON COLUMN pindah_unit.jabatan_id IS 'id jabatan terkait pada proses pindah (bila diisi)';
COMMENT ON COLUMN pindah_unit.keterangan IS 'Keterangan tambahan';
COMMENT ON COLUMN pindah_unit.nip IS 'NIP pegawai';
COMMENT ON COLUMN pindah_unit.no_sk_pindah IS 'Nomor SK pindah unit';
COMMENT ON COLUMN pindah_unit.sk_jabatan IS 'Nomor SK jabatan';
COMMENT ON COLUMN pindah_unit.sk_kp_terakhir IS 'Nomor SK kenaikan pangkat terakhir';
COMMENT ON COLUMN pindah_unit.sk_tunkin IS 'Nomor SK tunjangan kinerja';
COMMENT ON COLUMN pindah_unit.skp IS 'Sasaran kinerja pegawai';
COMMENT ON COLUMN pindah_unit.status_biro IS 'Status pengajuan pindah unit kerja di biro';
COMMENT ON COLUMN pindah_unit.status_satker IS 'Status pengajuan pindah unit kerja di satuan kerja';
COMMENT ON COLUMN pindah_unit.surat_permohonan_pindah IS 'Lokasi penyimpanan berkas permohonan pindah';
COMMENT ON COLUMN pindah_unit.surat_pernyataan_melepas IS 'Lokasi penyimpanan berkas pernyataan melepas pegawai';
COMMENT ON COLUMN pindah_unit.surat_pernyataan_menerima IS 'Lokasi penyimpanan berkas pernyataan menerima pegawai';
COMMENT ON COLUMN pindah_unit.tanggal_sk_pindah IS 'Tanggal SK pindah unit';
COMMENT ON COLUMN pindah_unit.tanggal_tmt_pindah IS 'TMT efektif perpindahan unit';
COMMENT ON COLUMN pindah_unit.unit_asal IS 'Unit kerja asal (rujuk unit_kerja)';
COMMENT ON COLUMN pindah_unit.unit_tujuan IS 'Unit kerja tujuan (rujuk unit_kerja)';
COMMENT ON COLUMN pindah_unit.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_agama
COMMENT ON COLUMN ref_agama.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_agama.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_agama.id IS 'id agama';
COMMENT ON COLUMN ref_agama.nama IS 'Nama agama';
COMMENT ON COLUMN ref_agama.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_golongan
COMMENT ON COLUMN ref_golongan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_golongan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_golongan.gol IS 'Golongan pada PNS';
COMMENT ON COLUMN ref_golongan.gol_pppk IS 'Golongan pada PPPK';
COMMENT ON COLUMN ref_golongan.id IS 'id golongan';
COMMENT ON COLUMN ref_golongan.nama IS 'Nama singkat golongan (mis. III/a)';
COMMENT ON COLUMN ref_golongan.nama_2 IS 'Nama singkat golongan lainnya (mis. 3a)';
COMMENT ON COLUMN ref_golongan.nama_pangkat IS 'Nama pangkat resmi';
COMMENT ON COLUMN ref_golongan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_instansi
COMMENT ON COLUMN ref_instansi.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_instansi.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_instansi.id IS 'id instansi';
COMMENT ON COLUMN ref_instansi.nama IS 'Nama instansi';
COMMENT ON COLUMN ref_instansi.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jabatan
COMMENT ON COLUMN ref_jabatan.bkn_id IS 'id pada sistem BKN';
COMMENT ON COLUMN ref_jabatan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jabatan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jabatan.id IS 'id jabatan'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_jabatan.jenis_jabatan IS 'Kode jenis jabatan';
COMMENT ON COLUMN ref_jabatan.kategori_jabatan IS 'Nama kategori jabatan';
COMMENT ON COLUMN ref_jabatan.kelas IS 'Kelas jabatan';
COMMENT ON COLUMN ref_jabatan.kode_bkn IS 'Kode pada sistem BKN'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_jabatan.kode_jabatan IS 'Kode unik jabatan'; -- TODO: ubah mapping kode_jabatan jadi id
COMMENT ON COLUMN ref_jabatan.nama_jabatan IS 'Nama jabatan';
COMMENT ON COLUMN ref_jabatan.nama_jabatan_bkn IS 'Nama jabatan pada sistem BKN';
COMMENT ON COLUMN ref_jabatan.nama_jabatan_full IS 'Nama jabatan lengkap';
COMMENT ON COLUMN ref_jabatan.no IS 'sama dengan id jabatan'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_jabatan.pensiun IS 'Usia pensiun jabatan terkait';
COMMENT ON COLUMN ref_jabatan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_diklat
COMMENT ON COLUMN ref_jenis_diklat.bkn_id IS 'id pada sistem BKN';
COMMENT ON COLUMN ref_jenis_diklat.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_diklat.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_diklat.id IS 'id jenis diklat';
COMMENT ON COLUMN ref_jenis_diklat.jenis_diklat IS 'Nama jenis diklat';
COMMENT ON COLUMN ref_jenis_diklat.kode IS 'kode jenis diklat';
COMMENT ON COLUMN ref_jenis_diklat.status IS 'status jenis diklat';
COMMENT ON COLUMN ref_jenis_diklat.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_hukuman
COMMENT ON COLUMN ref_jenis_hukuman.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_hukuman.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_hukuman.dikbud_hr_id IS 'id jenis hukuman disiplin pada Dikbud HR'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_jenis_hukuman.id IS 'id jenis hukuman disiplin';
COMMENT ON COLUMN ref_jenis_hukuman.nama IS 'Nama jenis hukuman disiplin';
COMMENT ON COLUMN ref_jenis_hukuman.nama_tingkat_hukuman IS 'Nama jenis tingkat hukuman disiplin';
COMMENT ON COLUMN ref_jenis_hukuman.tingkat_hukuman IS 'Nama pendek jenis hukuman disiplin, R: Ringan, S: Sedang, B: Berat';
COMMENT ON COLUMN ref_jenis_hukuman.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_jabatan
COMMENT ON COLUMN ref_jenis_jabatan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_jabatan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_jabatan.id IS 'id jenis jabatan';
COMMENT ON COLUMN ref_jenis_jabatan.nama IS 'Nama jenis jabatan';
COMMENT ON COLUMN ref_jenis_jabatan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_kawin
COMMENT ON COLUMN ref_jenis_kawin.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_kawin.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_kawin.id IS 'id status perkawinan';
COMMENT ON COLUMN ref_jenis_kawin.nama IS 'Nama status perkawinan';
COMMENT ON COLUMN ref_jenis_kawin.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_kp
COMMENT ON COLUMN ref_jenis_kp.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_kp.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_kp.dikbud_hr_id IS 'id jenis kenaikan pangkat pada Dikbud HR'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_jenis_kp.id IS 'id jenis kenaikan pangkat';
COMMENT ON COLUMN ref_jenis_kp.nama IS 'Nama jenis kenaikan pangkat';
COMMENT ON COLUMN ref_jenis_kp.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_pegawai
COMMENT ON COLUMN ref_jenis_pegawai.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_pegawai.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_pegawai.dikbud_hr_id IS 'id jenis pegawai pada Dikbud HR'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_jenis_pegawai.id IS 'id jenis pegawai';
COMMENT ON COLUMN ref_jenis_pegawai.nama IS 'Nama jenis pegawai';
COMMENT ON COLUMN ref_jenis_pegawai.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_penghargaan
COMMENT ON COLUMN ref_jenis_penghargaan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_penghargaan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_penghargaan.id IS 'id jenis penghargaan';
COMMENT ON COLUMN ref_jenis_penghargaan.nama IS 'Nama jenis penghargaan';
COMMENT ON COLUMN ref_jenis_penghargaan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_jenis_satker
COMMENT ON COLUMN ref_jenis_satker.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_jenis_satker.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_jenis_satker.id IS 'id jenis satuan kerja';
COMMENT ON COLUMN ref_jenis_satker.nama IS 'Nama jenis satuan kerja';
COMMENT ON COLUMN ref_jenis_satker.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_kedudukan_hukum
COMMENT ON COLUMN ref_kedudukan_hukum.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_kedudukan_hukum.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_kedudukan_hukum.dikbud_hr_id IS 'id kedudukan hukum pada Dikbud HR'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN ref_kedudukan_hukum.id IS 'id kedudukan hukum';
COMMENT ON COLUMN ref_kedudukan_hukum.is_pppk IS 'Apakah kedudukan hukum merupakan PPPK';
COMMENT ON COLUMN ref_kedudukan_hukum.nama IS 'Nama kedudukan hukum';
COMMENT ON COLUMN ref_kedudukan_hukum.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_kelas_jabatan
COMMENT ON COLUMN ref_kelas_jabatan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_kelas_jabatan.id IS 'id kelas jabatan';
COMMENT ON COLUMN ref_kelas_jabatan.kelas_jabatan IS 'Nama kelas jabatan';
COMMENT ON COLUMN ref_kelas_jabatan.tunjangan_kinerja IS 'Nilai tunjangan kinerja pada kelas jabatan terkait';
COMMENT ON COLUMN ref_kelas_jabatan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_kpkn
COMMENT ON COLUMN ref_kpkn.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_kpkn.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_kpkn.id IS 'id KPPN/KPKN';
COMMENT ON COLUMN ref_kpkn.nama IS 'Nama KPPN/KPKN';
COMMENT ON COLUMN ref_kpkn.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_lokasi
COMMENT ON COLUMN ref_lokasi.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_lokasi.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_lokasi.ibukota IS 'Nama ibu kota';
COMMENT ON COLUMN ref_lokasi.id IS 'id lokasi';
COMMENT ON COLUMN ref_lokasi.jenis IS 'Jenis lokasi, P: Provinsi, KA: Kota, KB: Kabupaten, KC: Kecamatan, DE: Desa, NG: Negara';
COMMENT ON COLUMN ref_lokasi.jenis_desa IS 'Jenis desa, K: Kelurahan, D: Desa';
COMMENT ON COLUMN ref_lokasi.jenis_kabupaten IS 'Jenis kabupaten, KOT: Kota, KBP: Kabupaten';
COMMENT ON COLUMN ref_lokasi.kanreg_id IS 'id kantor regional';
COMMENT ON COLUMN ref_lokasi.lokasi_id IS 'id lokasi wilayah administratif 1 tingkat di atasnya';
COMMENT ON COLUMN ref_lokasi.nama IS 'Nama lokasi';
COMMENT ON COLUMN ref_lokasi.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_pendidikan
COMMENT ON COLUMN ref_pendidikan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_pendidikan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_pendidikan.id IS 'id pendidikan (UUID)';
COMMENT ON COLUMN ref_pendidikan.nama IS 'Nama program/jurusan/pendidikan';
COMMENT ON COLUMN ref_pendidikan.tingkat_pendidikan_id IS 'Tingkat pendidikan (rujuk tingkat_pendidikan)';
COMMENT ON COLUMN ref_pendidikan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: ref_tingkat_pendidikan
COMMENT ON COLUMN ref_tingkat_pendidikan.abbreviation IS 'Singkatan tingkat pendidikan';
COMMENT ON COLUMN ref_tingkat_pendidikan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN ref_tingkat_pendidikan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN ref_tingkat_pendidikan.golongan_id IS 'Kaitan ke golongan awal/rujukan (jika relevan)';
COMMENT ON COLUMN ref_tingkat_pendidikan.id IS 'Kode/ID tingkat pendidikan';
COMMENT ON COLUMN ref_tingkat_pendidikan.nama IS 'Nama tingkat pendidikan (mis. S1, S2)';
COMMENT ON COLUMN ref_tingkat_pendidikan.golongan_awal_id IS 'id golongan awal';
COMMENT ON COLUMN ref_tingkat_pendidikan.tingkat IS 'Urutan/level numerik pendidikan';
COMMENT ON COLUMN ref_tingkat_pendidikan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_assesmen
COMMENT ON COLUMN riwayat_assesmen.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_assesmen.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_assesmen.file_upload IS 'Lokasi penyimpanan berkas asesmen';
COMMENT ON COLUMN riwayat_assesmen.file_upload_exists IS 'Penanda apakah berkas telah diunggah';
COMMENT ON COLUMN riwayat_assesmen.file_upload_fb_potensi IS 'Lokasi penyimpanan berkas umpan balik asesmen pada asesmen-pegawai.kemendikdasmen.go.id';
COMMENT ON COLUMN riwayat_assesmen.file_upload_fb_pt IS 'Lokasi penyimpanan berkas umpan balik asesmen pada asesmen-pegawai.kemendikdasmen.go.id'; -- TODO: hapus jika sudah tidak digunakan
COMMENT ON COLUMN riwayat_assesmen.file_upload_lengkap_pt IS 'Lokasi penyimpanan berkas lengkap hasil asesmen pada asesmen-pegawai.kemendikdasmen.go.id';
COMMENT ON COLUMN riwayat_assesmen.id IS 'id data asesmen';
COMMENT ON COLUMN riwayat_assesmen.nama_lengkap IS 'Nama lengkap pegawai yang diases';
COMMENT ON COLUMN riwayat_assesmen.nama_unor IS 'Nama unit organisasi pegawai yang diases';
COMMENT ON COLUMN riwayat_assesmen.nilai IS 'Hasil penilaian asesmen';
COMMENT ON COLUMN riwayat_assesmen.nilai_kinerja IS 'Hasil penilaian kinerja';
COMMENT ON COLUMN riwayat_assesmen.pns_id IS 'id PNS';
COMMENT ON COLUMN riwayat_assesmen.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_assesmen.posisi_id IS 'id posisi';
COMMENT ON COLUMN riwayat_assesmen.saran_pengembangan IS 'Saran pengembangan';
COMMENT ON COLUMN riwayat_assesmen.satker_id IS 'id satuan kerja';
COMMENT ON COLUMN riwayat_assesmen.tahun IS 'Tahun asesmen';
COMMENT ON COLUMN riwayat_assesmen.tahun_penilaian_id IS 'id tahun penilaian';
COMMENT ON COLUMN riwayat_assesmen.tahun_penilaian_title IS 'Judul tahun pada laporan hasil asesmen';
COMMENT ON COLUMN riwayat_assesmen.unit_org_id IS 'id unit organisasi';
COMMENT ON COLUMN riwayat_assesmen.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_diklat
COMMENT ON COLUMN riwayat_diklat.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_diklat.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_diklat.diklat_struktural_id IS 'id referensi diklat struktural';
COMMENT ON COLUMN riwayat_diklat.durasi_jam IS 'Durasi pelatihan dalam jam';
COMMENT ON COLUMN riwayat_diklat.file_base64 IS 'Berkas sertifikat diklat dalam format base64';
COMMENT ON COLUMN riwayat_diklat.id IS 'id riwayat diklat';
COMMENT ON COLUMN riwayat_diklat.institusi_penyelenggara IS 'Nama lembaga penyelenggara diklat';
COMMENT ON COLUMN riwayat_diklat.jenis_diklat IS 'Nama jenis diklat';
COMMENT ON COLUMN riwayat_diklat.jenis_diklat_id IS 'Jenis diklat (rujuk ref_jenis_diklat)';
COMMENT ON COLUMN riwayat_diklat.nama_diklat IS 'Nama diklat';
COMMENT ON COLUMN riwayat_diklat.nip_baru IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_diklat.no_sertifikat IS 'Nomor sertifikat diklat';
COMMENT ON COLUMN riwayat_diklat.pns_orang_id IS 'Referensi pegawai (rujuk pegawai.pns_id)';
COMMENT ON COLUMN riwayat_diklat.rumpun_diklat_id IS 'id rumpun diklat';
COMMENT ON COLUMN riwayat_diklat.rumpun_diklat_nama IS 'Nama rumpun diklat';
COMMENT ON COLUMN riwayat_diklat.siasn_id IS 'id referensi pada sistem BKN';
COMMENT ON COLUMN riwayat_diklat.sudah_kirim_siasn IS 'Penanda data sudah dikirim ke BKN';
COMMENT ON COLUMN riwayat_diklat.tahun_diklat IS 'Tahun pelaksanaan diklat';
COMMENT ON COLUMN riwayat_diklat.tanggal_mulai IS 'Tanggal mulai diklat';
COMMENT ON COLUMN riwayat_diklat.tanggal_selesai IS 'Tanggal selesai diklat';
COMMENT ON COLUMN riwayat_diklat.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_golongan
COMMENT ON COLUMN riwayat_golongan.arsip_id IS 'id referensi arsip';
COMMENT ON COLUMN riwayat_golongan.basic IS 'Penanda golongan basic';
COMMENT ON COLUMN riwayat_golongan.bkn_id IS 'id pada sistem BKN';
COMMENT ON COLUMN riwayat_golongan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_golongan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_golongan.file_base64 IS 'Berkas dalam format base64';
COMMENT ON COLUMN riwayat_golongan.golongan_asal IS 'Golongan asal pegawai';
COMMENT ON COLUMN riwayat_golongan.golongan_id IS 'id golongan pegawai';
COMMENT ON COLUMN riwayat_golongan.golongan_nama IS 'Nama golongan pegawai';
COMMENT ON COLUMN riwayat_golongan.id IS 'id riwayat golongan';
COMMENT ON COLUMN riwayat_golongan.jenis_kp IS 'Jenis kp';
COMMENT ON COLUMN riwayat_golongan.jenis_kp_id IS 'Jenis kp (rujuk ref_jenis_kp)';
COMMENT ON COLUMN riwayat_golongan.jenis_riwayat IS 'Jenis riwayat';
COMMENT ON COLUMN riwayat_golongan.jumlah_angka_kredit_tambahan IS 'Jumlah angka kredit tambahan';
COMMENT ON COLUMN riwayat_golongan.jumlah_angka_kredit_utama IS 'Jumlah angka kredit utama';
COMMENT ON COLUMN riwayat_golongan.kanreg IS 'Penanda apakah pegawai memiliki keterangan reguler';
COMMENT ON COLUMN riwayat_golongan.keterangan IS 'Keterangan golongan';
COMMENT ON COLUMN riwayat_golongan.keterangan_berkas IS 'Keterangan berkas';
COMMENT ON COLUMN riwayat_golongan.kode_jenis_kp IS 'Kode jenis kp';
COMMENT ON COLUMN riwayat_golongan.kpkn IS 'Keterangan kpkn';
COMMENT ON COLUMN riwayat_golongan.lpnk IS 'Keterangan lpnk';
COMMENT ON COLUMN riwayat_golongan.mk_golongan_bulan IS 'Bulan pemberian golongan';
COMMENT ON COLUMN riwayat_golongan.mk_golongan_tahun IS 'Tahun pemberian golongan';
COMMENT ON COLUMN riwayat_golongan.no_bkn IS 'Nomor BKN';
COMMENT ON COLUMN riwayat_golongan.pangkat_nama IS 'Nama pangkat';
COMMENT ON COLUMN riwayat_golongan.pangkat_terakhir IS 'Penanda apakah golongan merupakan golongan terakhir';
COMMENT ON COLUMN riwayat_golongan.pns_id IS 'Referensi pegawai (rujuk pegawai)';
COMMENT ON COLUMN riwayat_golongan.pns_nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_golongan.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_golongan.sk_nomor IS 'Nomor SK pemberian golongan';
COMMENT ON COLUMN riwayat_golongan.sk_tanggal IS 'Tanggal SK pemberian golongan';
COMMENT ON COLUMN riwayat_golongan.sk_type IS 'Jenis SK pemberian golongan';
COMMENT ON COLUMN riwayat_golongan.status_biro IS 'Status verifikasi di tingkat biro';
COMMENT ON COLUMN riwayat_golongan.status_satker IS 'Status verifikasi di tingkat satuan kerja';
COMMENT ON COLUMN riwayat_golongan.tanggal_bkn IS 'Tanggal BKN';
COMMENT ON COLUMN riwayat_golongan.tmt_golongan IS 'Tanggal mulai efektif golongan';
COMMENT ON COLUMN riwayat_golongan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_hukdis
COMMENT ON COLUMN riwayat_hukdis.bkn_id IS 'id pada sistem BKN';
COMMENT ON COLUMN riwayat_hukdis.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_hukdis.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_hukdis.file_base64 IS 'Berkas dalam format base64';
COMMENT ON COLUMN riwayat_hukdis.golongan_id IS 'id golongan pegawai';
COMMENT ON COLUMN riwayat_hukdis.id IS 'id riwayat hukuman disiplin';
COMMENT ON COLUMN riwayat_hukdis.jenis_hukuman_id IS 'id jenis hukuman (rujuk ref_jenis_hukuman)';
COMMENT ON COLUMN riwayat_hukdis.keterangan_berkas IS 'Keterangan berkas';
COMMENT ON COLUMN riwayat_hukdis.masa_bulan IS 'Masa hukuman dalam bulan';
COMMENT ON COLUMN riwayat_hukdis.masa_tahun IS 'Masa hukuman dalam tahun';
COMMENT ON COLUMN riwayat_hukdis.nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_hukdis.nama_golongan IS 'Nama golongan pegawai';
COMMENT ON COLUMN riwayat_hukdis.nama_jenis_hukuman IS 'Nama jenis hukuman';
COMMENT ON COLUMN riwayat_hukdis.no_pp IS 'Nomor PP pegawai';
COMMENT ON COLUMN riwayat_hukdis.no_sk_pembatalan IS 'Nomor SK pembatalan';
COMMENT ON COLUMN riwayat_hukdis.pns_id IS 'Referensi pegawai (rujuk pegawai.pns_id)';
COMMENT ON COLUMN riwayat_hukdis.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_hukdis.sk_nomor IS 'Nomor SK hukuman';
COMMENT ON COLUMN riwayat_hukdis.sk_tanggal IS 'Tanggal SK hukuman';
COMMENT ON COLUMN riwayat_hukdis.tanggal_akhir_hukuman IS 'Tanggal akhir masa hukuman';
COMMENT ON COLUMN riwayat_hukdis.tanggal_mulai_hukuman IS 'Tanggal mulai masa hukuman';
COMMENT ON COLUMN riwayat_hukdis.tanggal_sk_pembatalan IS 'Tanggal SK pembatalan';
COMMENT ON COLUMN riwayat_hukdis.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_jabatan
COMMENT ON COLUMN riwayat_jabatan.bkn_id IS 'id pada sistem BKN';
COMMENT ON COLUMN riwayat_jabatan.catatan IS 'Catatan atas riwayat jabatan';
COMMENT ON COLUMN riwayat_jabatan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_jabatan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_jabatan.eselon IS 'Nama eselon jabatan';
COMMENT ON COLUMN riwayat_jabatan.eselon_id IS 'id eselon jabatan';
COMMENT ON COLUMN riwayat_jabatan.eselon1 IS 'Unit eselon 1 terkait jabatan';
COMMENT ON COLUMN riwayat_jabatan.eselon2 IS 'Unit eselon 2 terkait jabatan';
COMMENT ON COLUMN riwayat_jabatan.eselon3 IS 'Unit eselon 3 terkait jabatan';
COMMENT ON COLUMN riwayat_jabatan.eselon4 IS 'Unit eselon 4 terkait jabatan';
COMMENT ON COLUMN riwayat_jabatan.id IS 'id riwayat jabatan';
COMMENT ON COLUMN riwayat_jabatan.is_active IS 'Penanda apakah jabatan masih aktif saat ini';
COMMENT ON COLUMN riwayat_jabatan.jabatan_id IS 'id jabatan (rujuk ref_jabatan)';
COMMENT ON COLUMN riwayat_jabatan.jabatan_id_bkn IS 'id jabatan pada sistem BKN';
COMMENT ON COLUMN riwayat_jabatan.jenis_jabatan IS 'Nama jenis jabatan';
COMMENT ON COLUMN riwayat_jabatan.jenis_jabatan_id IS 'id jenis jabatan (struktural/fungsional/dll)';
COMMENT ON COLUMN riwayat_jabatan.jenis_sk IS 'Kategori/jenis SK jabatan';
COMMENT ON COLUMN riwayat_jabatan.kelas_jabatan_id IS 'id kelas jabatan';
COMMENT ON COLUMN riwayat_jabatan.nama_jabatan IS 'Nama jabatan (teks)';
COMMENT ON COLUMN riwayat_jabatan.no_sk IS 'Nomor SK jabatan';
COMMENT ON COLUMN riwayat_jabatan.periode_jabatan_end_date IS 'Tanggal akhir periode jabatan';
COMMENT ON COLUMN riwayat_jabatan.periode_jabatan_start_date IS 'Tanggal mulai periode jabatan';
COMMENT ON COLUMN riwayat_jabatan.pns_id IS 'Referensi pegawai (rujuk pegawai.pns_id)';
COMMENT ON COLUMN riwayat_jabatan.pns_nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_jabatan.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_jabatan.satuan_kerja_id IS 'Satuan kerja terkait jabatan (rujuk unit_kerja)';
COMMENT ON COLUMN riwayat_jabatan.status_biro IS 'Status persetujuan biro kepegawaian';
COMMENT ON COLUMN riwayat_jabatan.status_plt IS 'Status pelaksana tugas (PLT)';
COMMENT ON COLUMN riwayat_jabatan.status_satker IS 'Status persetujuan satuan kerja';
COMMENT ON COLUMN riwayat_jabatan.tabel_mutasi_id IS 'Referensi ke tabel mutasi';
COMMENT ON COLUMN riwayat_jabatan.tanggal_sk IS 'Tanggal SK jabatan';
COMMENT ON COLUMN riwayat_jabatan.tmt_jabatan IS 'Tanggal mulai memangku jabatan';
COMMENT ON COLUMN riwayat_jabatan.tmt_pelantikan IS 'Tanggal mulai pelantikan';
COMMENT ON COLUMN riwayat_jabatan.unor IS 'Nama unit organisasi';
COMMENT ON COLUMN riwayat_jabatan.unor_id IS 'id unit organisasi saat jabatan (rujuk unit_kerja)';
COMMENT ON COLUMN riwayat_jabatan.unor_id_bkn IS 'id unit organisasi saat jabatan pada sistem BKN';
COMMENT ON COLUMN riwayat_jabatan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_kgb
COMMENT ON COLUMN riwayat_kgb.alasan IS 'Alasan kenaikan gaji berkala';
COMMENT ON COLUMN riwayat_kgb.birth_date IS 'Tanggal lahir pegawai';
COMMENT ON COLUMN riwayat_kgb.birth_place IS 'Tempat lahir pegawai';
COMMENT ON COLUMN riwayat_kgb.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_kgb.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_kgb.id IS 'id riwayat kenaikan gaji berkala';
COMMENT ON COLUMN riwayat_kgb.kantor_pembayaran IS 'Kantor yang melakukan pembayaran gaji';
COMMENT ON COLUMN riwayat_kgb.last_education IS 'Pendidikan terakhir';
COMMENT ON COLUMN riwayat_kgb.last_education_date IS 'Tanggal lulus pendidikan terakhir';
COMMENT ON COLUMN riwayat_kgb.mv_kgb_id IS 'Referensi ke materialized view KGB';
COMMENT ON COLUMN riwayat_kgb.n_gapok IS 'Nilai gaji pokok';
COMMENT ON COLUMN riwayat_kgb.n_gol_ruang IS 'Golongan ruang';
COMMENT ON COLUMN riwayat_kgb.n_gol_tmt IS 'Tanggal mulai golongan terkait';
COMMENT ON COLUMN riwayat_kgb.n_golongan_id IS 'Referensi ke golongan';
COMMENT ON COLUMN riwayat_kgb.n_jabatan_text IS 'Nama jabatan';
COMMENT ON COLUMN riwayat_kgb.n_masakerja_bln IS 'Masa kerja dalam bulan';
COMMENT ON COLUMN riwayat_kgb.n_masakerja_thn IS 'Masa kerja dalam tahun';
COMMENT ON COLUMN riwayat_kgb.n_tmt_jabatan IS 'Tanggal mulai memangku jabatan';
COMMENT ON COLUMN riwayat_kgb.no_sk IS 'Nomor SK kenaikan gaji berkala';
COMMENT ON COLUMN riwayat_kgb.pegawai_id IS 'Referensi ke data pegawai';
COMMENT ON COLUMN riwayat_kgb.pegawai_nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_kgb.pegawai_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_kgb.pejabat IS 'Pejabat yang menandatangani';
COMMENT ON COLUMN riwayat_kgb.ref IS 'Nomor referensi';
COMMENT ON COLUMN riwayat_kgb.tgl_sk IS 'Tanggal SK kenaikan gaji berkala';
COMMENT ON COLUMN riwayat_kgb.tmt_sk IS 'TMT SK kenaikan gaji berkala';
COMMENT ON COLUMN riwayat_kgb.unit_kerja_induk_id IS 'id unit kerja induk';
COMMENT ON COLUMN riwayat_kgb.unit_kerja_induk_text IS 'Nama unit kerja induk';
COMMENT ON COLUMN riwayat_kgb.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_kinerja
COMMENT ON COLUMN riwayat_kinerja.created_at IS 'Waktu pembuatan data';
COMMENT ON COLUMN riwayat_kinerja.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_kinerja.id IS 'id riwayat kinerja';
COMMENT ON COLUMN riwayat_kinerja.jabatan_penilai IS 'Jabatan pejabat penilai';
COMMENT ON COLUMN riwayat_kinerja.jabatan_penilai_realisasi IS 'Jabatan pejabat penilai hasil aktual';
COMMENT ON COLUMN riwayat_kinerja.nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_kinerja.nama_penilai IS 'Nama pejabat penilai';
COMMENT ON COLUMN riwayat_kinerja.nama_penilai_realisasi IS 'Nama pejabat penilai hasil aktual';
COMMENT ON COLUMN riwayat_kinerja.nip IS 'Nomor Induk Pegawai';
COMMENT ON COLUMN riwayat_kinerja.nip_penilai IS 'NIP pejabat penilai';
COMMENT ON COLUMN riwayat_kinerja.nip_penilai_realisasi IS 'NIP pejabat penilai hasil aktual';
COMMENT ON COLUMN riwayat_kinerja.predikat_kinerja IS 'Predikat penilaian kinerja';
COMMENT ON COLUMN riwayat_kinerja.rating_hasil_kerja IS 'Nilai hasil kerja';
COMMENT ON COLUMN riwayat_kinerja.rating_perilaku_kerja IS 'Nilai perilaku kerja';
COMMENT ON COLUMN riwayat_kinerja.ref IS 'Nomor referensi';
COMMENT ON COLUMN riwayat_kinerja.tahun IS 'Tahun penilaian';
COMMENT ON COLUMN riwayat_kinerja.unit_kerja_penilai IS 'Unit kerja pejabat penilai';
COMMENT ON COLUMN riwayat_kinerja.unit_kerja_penilai_realisasi IS 'Unit kerja pejabat penilai hasil aktual';
COMMENT ON COLUMN riwayat_kinerja.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_kursus
COMMENT ON COLUMN riwayat_kursus.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_kursus.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_kursus.id IS 'id riwayat kursus';
COMMENT ON COLUMN riwayat_kursus.instansi IS 'Nama instansi penyelenggara';
COMMENT ON COLUMN riwayat_kursus.institusi_penyelenggara IS 'Nama lembaga yang menyelenggarakan';
COMMENT ON COLUMN riwayat_kursus.jenis_kursus IS 'Jenis kursus';
COMMENT ON COLUMN riwayat_kursus.lama_kursus IS 'Durasi pelaksanaan';
COMMENT ON COLUMN riwayat_kursus.nama_kursus IS 'Nama kursus';
COMMENT ON COLUMN riwayat_kursus.no_sertifikat IS 'Nomor sertifikat yang diterbitkan';
COMMENT ON COLUMN riwayat_kursus.pns_id IS 'id PNS';
COMMENT ON COLUMN riwayat_kursus.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_kursus.siasn_id IS 'id pada sistem BKN';
COMMENT ON COLUMN riwayat_kursus.tanggal_kursus IS 'Tanggal pelaksanaan';
COMMENT ON COLUMN riwayat_kursus.tipe_kursus IS 'Tipe kursus';
COMMENT ON COLUMN riwayat_kursus.updated_at IS 'Waktu terakhir data diperbarui';

-- Comments for table: riwayat_nine_box
COMMENT ON COLUMN riwayat_nine_box.created_at IS 'Waktu pembuatan data';
COMMENT ON COLUMN riwayat_nine_box.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_nine_box.id IS 'id data nine box';
COMMENT ON COLUMN riwayat_nine_box.kelas_jabatan IS 'Kelas jabatan pegawai';
COMMENT ON COLUMN riwayat_nine_box.kesimpulan IS 'Kesimpulan hasil penilaian';
COMMENT ON COLUMN riwayat_nine_box.nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_nine_box.nama_jabatan IS 'Nama jabatan pegawai';
COMMENT ON COLUMN riwayat_nine_box.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_nine_box.tahun IS 'Tahun penilaian';
COMMENT ON COLUMN riwayat_nine_box.updated_at IS 'Waktu terakhir pembaruan data';

-- Comments for table: riwayat_pendidikan
COMMENT ON COLUMN riwayat_pendidikan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_pendidikan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_pendidikan.diakui_bkn IS 'Penanda pengakuan oleh BKN';
COMMENT ON COLUMN riwayat_pendidikan.file_base64 IS 'File berkas dalam format base64';
COMMENT ON COLUMN riwayat_pendidikan.gelar_belakang IS 'Gelar belakang terkait pendidikan';
COMMENT ON COLUMN riwayat_pendidikan.gelar_depan IS 'Gelar depan terkait pendidikan';
COMMENT ON COLUMN riwayat_pendidikan.id IS 'id riwayat pendidikan pegawai';
COMMENT ON COLUMN riwayat_pendidikan.keterangan_berkas IS 'Keterangan berkas';
COMMENT ON COLUMN riwayat_pendidikan.nama_sekolah IS 'Nama institusi pendidikan';
COMMENT ON COLUMN riwayat_pendidikan.negara_sekolah IS 'Negara tempat sekolah/pendidikan';
COMMENT ON COLUMN riwayat_pendidikan.nip IS 'Nomor Induk Pegawai';
COMMENT ON COLUMN riwayat_pendidikan.no_ijazah IS 'Nomor ijazah';
COMMENT ON COLUMN riwayat_pendidikan.pendidikan_id IS 'Program pendidikan (rujuk pendidikan)';
COMMENT ON COLUMN riwayat_pendidikan.pendidikan_id_3 IS 'id pendidikan versi 3';
COMMENT ON COLUMN riwayat_pendidikan.pendidikan_pertama IS 'Penanda pendidikan pertama';
COMMENT ON COLUMN riwayat_pendidikan.pendidikan_terakhir IS 'Penanda pendidikan terakhir';
COMMENT ON COLUMN riwayat_pendidikan.pns_id IS 'Referensi pegawai (rujuk pegawai.pns_id)';
COMMENT ON COLUMN riwayat_pendidikan.pns_id_3 IS 'id PNS versi 3';
COMMENT ON COLUMN riwayat_pendidikan.status_biro IS 'Status verifikasi di tingkat biro';
COMMENT ON COLUMN riwayat_pendidikan.status_satker IS 'Status verifikasi di tingkat satuan kerja';
COMMENT ON COLUMN riwayat_pendidikan.tahun_lulus IS 'Tahun kelulusan';
COMMENT ON COLUMN riwayat_pendidikan.tanggal_lulus IS 'Tanggal kelulusan';
COMMENT ON COLUMN riwayat_pendidikan.tingkat_pendidikan_id IS 'Tingkat pendidikan (rujuk tingkat_pendidikan)';
COMMENT ON COLUMN riwayat_pendidikan.tugas_belajar IS 'Tugas belajar';
COMMENT ON COLUMN riwayat_pendidikan.updated_at IS 'Waktu terakhir pembaruan';

-- Comments for table: riwayat_penghargaan_umum
COMMENT ON COLUMN riwayat_penghargaan_umum.created_at IS 'Waktu pembuatan data';
COMMENT ON COLUMN riwayat_penghargaan_umum.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_penghargaan_umum.deskripsi_penghargaan IS 'Deskripsi lengkap penghargaan yang diterima';
COMMENT ON COLUMN riwayat_penghargaan_umum.exist IS 'Penanda keberadaan data';
COMMENT ON COLUMN riwayat_penghargaan_umum.file_base64 IS 'Berkas bukti penghargaan dalam format base64';
COMMENT ON COLUMN riwayat_penghargaan_umum.id IS 'ID unik data penghargaan';
COMMENT ON COLUMN riwayat_penghargaan_umum.jenis_penghargaan_id IS 'Jenis atau kategori penghargaan (rujuk ref_jenis_penghargaan)';
COMMENT ON COLUMN riwayat_penghargaan_umum.nama_penghargaan IS 'Nama penghargaan yang diterima';
COMMENT ON COLUMN riwayat_penghargaan_umum.nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_penghargaan_umum.tanggal_penghargaan IS 'Tanggal penerimaan penghargaan';
COMMENT ON COLUMN riwayat_penghargaan_umum.updated_at IS 'Waktu terakhir pembaruan data';

-- Comments for table: riwayat_penugasan
COMMENT ON COLUMN riwayat_penugasan.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_penugasan.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_penugasan.deskripsi_jabatan IS 'Deskripsi jabatan';
COMMENT ON COLUMN riwayat_penugasan.file_base64 IS 'Berkas bukti penugasan dalam format base64';
COMMENT ON COLUMN riwayat_penugasan.id IS 'ID unik data penugasan';
COMMENT ON COLUMN riwayat_penugasan.is_menjabat IS 'Penanda apakah pegawai masih menjabat';
COMMENT ON COLUMN riwayat_penugasan.nama_jabatan IS 'Nama jabatan';
COMMENT ON COLUMN riwayat_penugasan.nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_penugasan.tanggal_mulai IS 'Tanggal mulai penugasan';
COMMENT ON COLUMN riwayat_penugasan.tanggal_selesai IS 'Tanggal selesai penugasan';
COMMENT ON COLUMN riwayat_penugasan.tipe_jabatan IS 'Tipe jabatan';
COMMENT ON COLUMN riwayat_penugasan.updated_at IS 'Waktu terakhir pembaruan data';

-- Comments for table: riwayat_pindah_unit_kerja
COMMENT ON COLUMN riwayat_pindah_unit_kerja.asal_id IS 'id unit kerja asal';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.asal_nama IS 'Nama unit kerja asal';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.file_base64 IS 'Berkas bukti pindah unit kerja dalam format base64';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.id IS 'id unik data pindah unit kerja';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.instansi_id IS 'id instansi tujuan';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.keterangan_berkas IS 'Keterangan berkas';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.nama_instansi IS 'Nama instansi tujuan';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.nama_satuan_kerja IS 'Nama satuan kerja tujuan';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.nama_unor_baru IS 'Nama unit organisasi baru';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.pns_id IS 'id pegawai';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.pns_nama IS 'Nama pegawai';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.pns_nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.satuan_kerja_id IS 'id satuan kerja tujuan';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.sk_nomor IS 'Nomor SK pindah unit kerja';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.sk_tanggal IS 'Tanggal SK pindah unit kerja';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.unor_id_baru IS 'id unit organisasi baru';
COMMENT ON COLUMN riwayat_pindah_unit_kerja.updated_at IS 'Waktu terakhir pembaruan data';

-- Comments for table: riwayat_sertifikasi
COMMENT ON COLUMN riwayat_sertifikasi.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_sertifikasi.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_sertifikasi.deskripsi IS 'Deskripsi sertifikasi';
COMMENT ON COLUMN riwayat_sertifikasi.file_base64 IS 'Berkas bukti sertifikasi dalam format base64';
COMMENT ON COLUMN riwayat_sertifikasi.id IS 'id unik data sertifikasi';
COMMENT ON COLUMN riwayat_sertifikasi.nama_sertifikasi IS 'Nama sertifikasi';
COMMENT ON COLUMN riwayat_sertifikasi.nip IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_sertifikasi.tahun IS 'Tahun sertifikasi';
COMMENT ON COLUMN riwayat_sertifikasi.updated_at IS 'Waktu terakhir pembaruan data';

-- Comments for table: riwayat_ujikom
COMMENT ON COLUMN riwayat_ujikom.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN riwayat_ujikom.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN riwayat_ujikom.exist IS 'Status apakah uji kompetensi sudah ada';
COMMENT ON COLUMN riwayat_ujikom.id IS 'id unik data uji kompetensi';
COMMENT ON COLUMN riwayat_ujikom.jenis_ujikom IS 'Jenis uji kompetensi';
COMMENT ON COLUMN riwayat_ujikom.link_sertifikat IS 'Link sertifikat';
COMMENT ON COLUMN riwayat_ujikom.nip_baru IS 'NIP pegawai';
COMMENT ON COLUMN riwayat_ujikom.tahun IS 'Tahun uji kompetensi';
COMMENT ON COLUMN riwayat_ujikom.updated_at IS 'Waktu terakhir pembaruan data';

-- Comments for table: unit_kerja
COMMENT ON COLUMN unit_kerja.abbreviation IS 'Singkatan unit organisasi';
COMMENT ON COLUMN unit_kerja.aktif IS 'Status keaktifan unit';
COMMENT ON COLUMN unit_kerja.cepat_kode IS 'Kode cepat untuk pencarian unit kerja';
COMMENT ON COLUMN unit_kerja.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN unit_kerja.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN unit_kerja.diatasan_id IS 'Unit atasan langsung (self-reference ke unit_kerja)';
COMMENT ON COLUMN unit_kerja.eselon_1 IS 'Kode eselon 1 unit kerja';
COMMENT ON COLUMN unit_kerja.eselon_2 IS 'Kode eselon 2 unit kerja';
COMMENT ON COLUMN unit_kerja.eselon_3 IS 'Kode eselon 3 unit kerja';
COMMENT ON COLUMN unit_kerja.eselon_4 IS 'Kode eselon 4 unit kerja';
COMMENT ON COLUMN unit_kerja.eselon_id IS 'id eselon unit (bila berlaku)';
COMMENT ON COLUMN unit_kerja.eselon_nama IS 'Nama eselon unit kerja';
COMMENT ON COLUMN unit_kerja.expired_date IS 'Tanggal kedaluwarsa unit kerja';
COMMENT ON COLUMN unit_kerja.id IS 'id unit organisasi (UUID)';
COMMENT ON COLUMN unit_kerja.instansi_id IS 'id instansi pemilik unit (rujuk ref_instansi)';
COMMENT ON COLUMN unit_kerja.is_satker IS 'Penanda apakah unit merupakan Satuan Kerja';
COMMENT ON COLUMN unit_kerja.jabatan_id IS 'ID jabatan yang terkait dengan unit kerja';
COMMENT ON COLUMN unit_kerja.jenis_satker IS 'Jenis satuan kerja';
COMMENT ON COLUMN unit_kerja.jenis_unor_id IS 'Jenis unit organisasi (bila digunakan)';
COMMENT ON COLUMN unit_kerja.jumlah_ideal_staff IS 'Jumlah ideal staf dalam unit kerja';
COMMENT ON COLUMN unit_kerja.keterangan IS 'Keterangan tambahan untuk unit kerja';
COMMENT ON COLUMN unit_kerja.kode_internal IS 'Kode internal unit organisasi';
COMMENT ON COLUMN unit_kerja.nama_jabatan IS 'Nama jabatan dalam unit kerja';
COMMENT ON COLUMN unit_kerja.nama_pejabat IS 'Nama pejabat yang menjabat';
COMMENT ON COLUMN unit_kerja.nama_unor IS 'Nama unit organisasi';
COMMENT ON COLUMN unit_kerja.no IS 'Nomor urut unit kerja';
COMMENT ON COLUMN unit_kerja.order IS 'Urutan tampilan unit kerja';
COMMENT ON COLUMN unit_kerja.pemimpin_pns_id IS 'ID PNS yang memimpin unit kerja';
COMMENT ON COLUMN unit_kerja.peraturan IS 'Peraturan yang mendasari unit kerja';
COMMENT ON COLUMN unit_kerja.remark IS 'Catatan tambahan untuk unit kerja';
COMMENT ON COLUMN unit_kerja.unor_induk IS 'Unit organisasi induk';
COMMENT ON COLUMN unit_kerja.unor_induk_penyetaraan IS 'Penyetaraan unit organisasi induk';
COMMENT ON COLUMN unit_kerja.updated_at IS 'Waktu terakhir data diperbarui';
COMMENT ON COLUMN unit_kerja.waktu IS 'Waktu pencatatan data unit kerja';

-- Comments for table: update_mandiri
COMMENT ON COLUMN update_mandiri.created_at IS 'Waktu perekaman data';
COMMENT ON COLUMN update_mandiri.dari IS 'Nilai asli kolom sebelum perubahan';
COMMENT ON COLUMN update_mandiri.deleted_at IS 'Waktu penghapusan data';
COMMENT ON COLUMN update_mandiri.id IS 'id unik data update mandiri';
COMMENT ON COLUMN update_mandiri.kolom IS 'Kolom yang diperbarui';
COMMENT ON COLUMN update_mandiri.level_update IS 'Level/update yang dilakukan';
COMMENT ON COLUMN update_mandiri.nama_kolom IS 'Nama kolom yang diperbarui';
COMMENT ON COLUMN update_mandiri.perubahan IS 'Nilai perubahan kolom';
COMMENT ON COLUMN update_mandiri.pns_id IS 'id pegawai';
COMMENT ON COLUMN update_mandiri.status IS 'Status verifikasi diperbarui';
COMMENT ON COLUMN update_mandiri.tabel_id IS 'id tabel yang diperbarui';
COMMENT ON COLUMN update_mandiri.updated_at IS 'Waktu terakhir pembaruan data';
COMMENT ON COLUMN update_mandiri.updated_by IS 'id user yang melakukan pembaruan';
COMMENT ON COLUMN update_mandiri.verifikasi_by IS 'id user yang melakukan verifikasi';
COMMENT ON COLUMN update_mandiri.verifikasi_tgl IS 'Tanggal verifikasi';
