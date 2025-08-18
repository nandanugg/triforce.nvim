create table kepegawaian.rwt_sertifikasi(
    id bigserial not null,
    nip varchar null,
    tahun int8 null,
    nama_sertifikasi varchar null,
    "base64" text null,
    createddate timestamp default now() null,
    deskripsi text null,
    constraint rwt_sertifikasi_pkey primary key (id)
);
