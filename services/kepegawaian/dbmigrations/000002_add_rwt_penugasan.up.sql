CREATE TABLE riwayat_penugasan (
    id serial PRIMARY KEY,
    tipe_jabatan varchar(200),
    deskripsi_jabatan varchar(200),
    tanggal_mulai date,
    tanggal_selesai date,
    file_base64 text,
    nip varchar(20),
    nama_jabatan varchar(200),
    is_menjabat bool,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz
);
