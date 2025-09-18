create table pemberitahuan(
    id bigserial primary key,
    judul_berita text not null,
    deskripsi_berita text not null,
    status text not null,
    updated_by text not null,
    updated_at date not null
);

create index pemberitahuan_judul_idx on pemberitahuan(judul_berita);

create index pemberitahuan_deskripsi_idx on pemberitahuan(deskripsi_berita);
