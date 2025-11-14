begin;

create unique index role_nama_unique_idx on role (lower(nama)) where deleted_at is null;

commit;
