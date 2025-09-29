-- name: GetDataPribadi :one
select
  pegawai.pns_id,
  pegawai.nama,
  pegawai.gelar_depan,
  pegawai.gelar_belakang,
  case
    when pegawai.jenis_jabatan_id = @jenis_jabatan_struktural::int and pegawai.pns_id = unit_kerja.pemimpin_pns_id then unit_kerja.nama_jabatan
    else ref_jabatan_aktual.nama_jabatan
  end as jabatan_aktual,
  case
    when pegawai.jenis_jabatan_id = @jenis_jabatan_struktural::int and pegawai.pns_id = unit_kerja.pemimpin_pns_id then 'Struktural'
    else ref_jenis_jabatan_aktual.nama
  end as jenis_jabatan_aktual,
  pegawai.tmt_jabatan,
  pegawai.unor_id,
  pegawai.nip_baru as nip,
  pegawai.nik,
  pegawai.kk,
  pegawai.jenis_kelamin,
  coalesce(ref_lokasi_tempat_lahir.nama, pegawai.tempat_lahir) as tempat_lahir,
  pegawai.tanggal_lahir,
  ref_tingkat_pendidikan.nama as tingkat_pendidikan,
  ref_pendidikan.nama as pendidikan,
  ref_jenis_kawin.nama as jenis_kawin,
  ref_agama.nama as agama,
  pegawai.email_dikbud,
  pegawai.email,
  pegawai.alamat,
  pegawai.no_hp,
  pegawai.no_darurat,
  ref_jenis_pegawai.nama as jenis_pegawai,
  case
    when pegawai.status_pegawai = 1 and (pegawai.terminated_date is null or pegawai.terminated_date > current_date)
      then (coalesce(extract(year from age(current_date, pegawai.tmt_cpns)), 0) + coalesce(pegawai.mk_tahun_swasta, 0) +
        floor((coalesce(extract(month from age(current_date, pegawai.tmt_cpns)), 0) + coalesce(pegawai.mk_bulan_swasta, 0)) / 12)) || ' Tahun ' ||
        ((coalesce(extract(month from age(current_date, pegawai.tmt_cpns)), 0) + coalesce(pegawai.mk_bulan_swasta, 0)) % 12) || ' Bulan'
  end as masa_kerja_keseluruhan,
  pegawai.masa_kerja as masa_kerja_golongan,
  ref_jabatan.nama_jabatan as jabatan,
  ref_jenis_jabatan.nama as jenis_jabatan,
  ref_kelas_jabatan.kelas_jabatan,
  coalesce(ref_lokasi_lokasi_kerja.nama, pegawai.lokasi_kerja) as lokasi_kerja,
  ref_golongan_akhir.nama_pangkat as pangkat_akhir,
  case when ref_kedudukan_hukum.is_pppk then ref_golongan_awal.gol_pppk else ref_golongan_awal.nama end as golongan_awal,
  case when ref_kedudukan_hukum.is_pppk then ref_golongan_akhir.gol_pppk else ref_golongan_akhir.nama end as golongan_akhir,
  pegawai.tmt_golongan,
  pegawai.tmt_cpns as tmt_asn,
  pegawai.no_sk_cpns as no_sk_asn,
  ref_kedudukan_hukum.is_pppk,
  ref_kedudukan_hukum.nama as status_asn,
  pegawai.status_cpns_pns as status_pns,
  pegawai.tmt_pns,
  pegawai.kartu_pegawai,
  pegawai.no_surat_dokter,
  pegawai.tanggal_surat_dokter,
  pegawai.no_bebas_narkoba,
  pegawai.tanggal_bebas_narkoba,
  pegawai.no_catatan_polisi,
  pegawai.tanggal_catatan_polisi,
  pegawai.akte_kelahiran,
  pegawai.bpjs,
  pegawai.npwp,
  pegawai.tanggal_npwp,
  pegawai.no_taspen,
  pegawai.foto
from pegawai
left join ref_lokasi ref_lokasi_tempat_lahir
	on ref_lokasi_tempat_lahir.id = pegawai.tempat_lahir_id and ref_lokasi_tempat_lahir.deleted_at is null
left join ref_agama
	on ref_agama.id = pegawai.agama_id and ref_agama.deleted_at is null
left join ref_tingkat_pendidikan
	on ref_tingkat_pendidikan.id = pegawai.tingkat_pendidikan_id and ref_tingkat_pendidikan.deleted_at is null
left join ref_pendidikan
	on ref_pendidikan.id = pegawai.pendidikan_id and ref_pendidikan.deleted_at is null
left join ref_lokasi ref_lokasi_lokasi_kerja
	on ref_lokasi_lokasi_kerja.id = pegawai.lokasi_kerja_id and ref_lokasi_lokasi_kerja.deleted_at is null
left join ref_jenis_pegawai
	on ref_jenis_pegawai.id = pegawai.jenis_pegawai_id and ref_jenis_pegawai.deleted_at is null
left join ref_golongan ref_golongan_akhir
	on ref_golongan_akhir.id = pegawai.gol_id and ref_golongan_akhir.deleted_at is null
left join ref_golongan ref_golongan_awal
	on ref_golongan_awal.id = pegawai.gol_awal_id and ref_golongan_awal.deleted_at is null
left join ref_kedudukan_hukum
	on ref_kedudukan_hukum.id = pegawai.kedudukan_hukum_id and ref_kedudukan_hukum.deleted_at is null
left join ref_jenis_kawin
	on ref_jenis_kawin.id = pegawai.jenis_kawin_id and ref_jenis_kawin.deleted_at is null
left join unit_kerja
	on unit_kerja.id = pegawai.unor_id and unit_kerja.deleted_at is null
left join ref_jabatan
	on ref_jabatan.kode_jabatan = pegawai.jabatan_instansi_id and ref_jabatan.deleted_at is null
left join ref_jenis_jabatan
	on ref_jenis_jabatan.id = ref_jabatan.jenis_jabatan and ref_jenis_jabatan.deleted_at is null
left join ref_kelas_jabatan
	on ref_kelas_jabatan.id = ref_jabatan.kelas
left join ref_jabatan ref_jabatan_aktual
	on ref_jabatan_aktual.kode_jabatan = pegawai.jabatan_instansi_real_id and ref_jabatan_aktual.deleted_at is null
left join ref_jenis_jabatan ref_jenis_jabatan_aktual
	on ref_jenis_jabatan_aktual.id = ref_jabatan_aktual.jenis_jabatan and ref_jenis_jabatan_aktual.deleted_at is null
where pegawai.nip_baru = $1 and pegawai.deleted_at is null;
