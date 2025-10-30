begin;

select setval('pasangan_id_seq', coalesce(max(id), 0) + 1, false) from pasangan;
select setval('anak_id_seq', coalesce(max(id), 0) + 1, false) from anak;
select setval('orang_tua_id_seq', coalesce(max(id), 0) + 1, false) from orang_tua;

alter table pasangan
  add column nik varchar(20),
  add column agama_id int2 references ref_agama(id);

comment on column pasangan.status is 'Referensi ke ref_jenis_kawin.id';

alter table anak
  add column nik varchar(20),
  add column agama_id int2 references ref_agama(id),
  add column jenis_kawin_id int2 references ref_jenis_kawin(id),
  add column status_sekolah int2,
  add column anak_ke int2;

comment on column anak.status_sekolah is 'Status Sekolah / Kerja, 1: Masih Sekolah, 2: Sudah Bekerja';

create index pegawai_nip_baru_idx on pegawai(nip_baru);
create index pasangan_pns_id_idx on pasangan(pns_id);
create index orang_tua_pns_id_idx on orang_tua(pns_id);
create index anak_pns_id_idx on anak(pns_id);

commit;
