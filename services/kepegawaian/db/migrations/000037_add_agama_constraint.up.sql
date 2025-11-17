begin;

create unique index agama_nama_unique_idx on ref_agama (lower(nama)) where deleted_at is null;

commit;
