-- name: GetDataPribadi :one
select
  p.nama,
  p.gelar_depan,
  p.gelar_belakang,
  case
    when p.jenis_jabatan_id = @jenis_jabatan_struktural::int and p.pns_id = uk.pemimpin_pns_id then uk.nama_jabatan
    else rjf.nama_jabatan
  end as jabatan_aktual,
  case
    when p.jenis_jabatan_id = @jenis_jabatan_struktural::int and p.pns_id = uk.pemimpin_pns_id then 'Struktural'
    else rjj.nama
  end as jenis_jabatan_aktual,
  p.tmt_jabatan,
  p.unor_id,
  p.nip_baru as nip,
  p.nik,
  p.kk,
  p.jenis_kelamin,
  coalesce(rltl.nama, p.tempat_lahir) as tempat_lahir,
  p.tgl_lahir,
  rtp.nama as tingkat_pendidikan,
  rp.nama as pendidikan,
  rjk.nama as jenis_kawin,
  ra.nama as agama,
  p.email_dikbud,
  p.email,
  p.alamat,
  p.no_hp,
  p.no_darurat,
  rjp.nama as jenis_pegawai,
  case
    when p.status_pegawai = 1 and (p.terminated_date is null or p.terminated_date > current_date)
      then (coalesce(extract(year from age(current_date, p.tmt_cpns)), 0) + coalesce(p.mk_tahun_swasta, 0) +
        floor((coalesce(extract(month from age(current_date, p.tmt_cpns)), 0) + coalesce(p.mk_bulan_swasta, 0)) / 12)) || ' Tahun ' ||
        ((coalesce(extract(month from age(current_date, p.tmt_cpns)), 0) + coalesce(p.mk_bulan_swasta, 0)) % 12) || ' Bulan'
  end as masa_kerja_keseluruhan,
  p.masa_kerja as masa_kerja_golongan,
  rj.nama_jabatan as jabatan,
  rkj.kelas_jabatan,
  coalesce(rlk.nama, p.lokasi_kerja) as lokasi_kerja,
  rge.nama_pangkat as pangkat_akhir,
  case when rkh.is_pppk then rgs.gol_pppk else rgs.nama end as golongan_awal,
  case when rkh.is_pppk then rge.gol_pppk else rge.nama end as golongan_akhir,
  p.tmt_golongan,
  p.tmt_cpns as tmt_asn,
  p.no_sk_cpns as no_sk_asn,
  case when rkh.is_pppk then rkh.nama else p.status_cpns_pns end as status_asn,
  case when rkh.is_pppk then null else p.tmt_pns end as tmt_pns,
  p.kartu_pegawai,
  p.no_surat_dokter,
  p.tgl_surat_dokter,
  p.no_bebas_narkoba,
  p.tgl_bebas_narkoba,
  p.no_catatan_polisi,
  p.tgl_catatan_polisi,
  p.akte_kelahiran,
  p.bpjs,
  p.npwp,
  p.tgl_npwp,
  p.no_taspen
from pegawai p
left join ref_lokasi rltl on rltl.id = p.tempat_lahir_id and rltl.deleted_at is null
left join ref_agama ra on ra.id = p.agama_id and ra.deleted_at is null
left join ref_tingkat_pendidikan rtp on rtp.id = p.tingkat_pendidikan_id and rtp.deleted_at is null
left join ref_pendidikan rp on rp.id = p.pendidikan_id and rp.deleted_at is null
left join ref_lokasi rlk on rlk.id = p.lokasi_kerja_id and rlk.deleted_at is null
left join ref_jenis_pegawai rjp on rjp.id = p.jenis_pegawai_id and rjp.deleted_at is null
left join ref_golongan rge on rge.id = p.gol_id and rge.deleted_at is null
left join ref_golongan rgs on rgs.id = p.gol_awal_id and rgs.deleted_at is null
left join ref_kedudukan_hukum rkh on rkh.id = p.kedudukan_hukum_id and rkh.deleted_at is null
left join ref_jenis_kawin rjk on rjk.id = p.jenis_kawin_id and rjk.deleted_at is null
left join unit_kerja uk on uk.id = p.unor_id and uk.deleted_at is null
left join ref_jabatan rj on rj.kode_jabatan = p.jabatan_instansi_id and rj.deleted_at is null
left join ref_kelas_jabatan rkj on rkj.id = rj.kelas
left join ref_jabatan rjf on rjf.kode_jabatan = p.jabatan_instansi_real_id and rjf.deleted_at is null
left join ref_jenis_jabatan rjj on rjj.id = rjf.jenis_jabatan and rjj.deleted_at is null
where p.nip_baru = $1 and p.deleted_at is null;
