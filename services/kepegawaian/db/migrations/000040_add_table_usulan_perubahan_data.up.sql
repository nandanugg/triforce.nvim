begin;

create table usulan_perubahan_data (
  id bigserial primary key,
  nip varchar(20) not null,
  jenis_data varchar(100) not null,
  data_id text,
  perubahan_data jsonb not null,
  action varchar(10) not null,
  status varchar(50) not null default 'Diusulkan',
  catatan varchar(200),
  read_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

create index usulan_perubahan_data_nip_idx on usulan_perubahan_data (nip);
create index usulan_perubahan_data_jenis_data_idx on usulan_perubahan_data (jenis_data);
create index usulan_perubahan_data_status_idx on usulan_perubahan_data (status);
create unique index usulan_perubahan_data_jenis_data_data_id_diusulkan_unique_idx on usulan_perubahan_data (jenis_data, data_id)
  where status = 'Diusulkan' and data_id is not null and deleted_at is null;

commit;
