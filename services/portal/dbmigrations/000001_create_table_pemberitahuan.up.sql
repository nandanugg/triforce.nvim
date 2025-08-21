create table portal.pemberitahuan(
    id bigserial primary key,
    judul_berita text not null,
    deskripsi_berita text not null,
    status text not null,
    updated_by text not null,
    updated_at date not null
);

create index pemberitahuan_judul_idx on portal.pemberitahuan(judul_berita);

create index pemberitahuan_deskripsi_idx on portal.pemberitahuan(deskripsi_berita);
