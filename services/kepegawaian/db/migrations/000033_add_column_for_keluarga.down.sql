begin;

alter table pasangan
  drop column nik,
  drop column agama_id;

comment on column pasangan.status is 'Status hubungan saat ini, 1: menikah, 2: cerai, 3: janda/duda';

alter table anak
  drop column nik,
  drop column agama_id,
  drop column jenis_kawin_id,
  drop column status_sekolah,
  drop column anak_ke;

drop index pegawai_nip_baru_idx;
drop index pasangan_pns_id_idx;
drop index orang_tua_pns_id_idx;
drop index anak_pns_id_idx;

commit;
