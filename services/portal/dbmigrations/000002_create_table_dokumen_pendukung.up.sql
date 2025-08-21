create table portal.dokumen_pendukung(
    id bigserial primary key,
    nama_tombol text not null,
    nama_halaman text not null,
    file bytea,
    updated_by text,
    updated_at date
);
