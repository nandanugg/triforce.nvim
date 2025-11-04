BEGIN;

CREATE TABLE pegawai_ttd (
    pns_id VARCHAR(36) PRIMARY KEY,
    nip VARCHAR(20) NOT NULL,
    base64ttd TEXT NOT NULL,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create index pegawai_ttd_nip_idx on pegawai_ttd(nip);

alter table surat_keputusan 
    ALTER COLUMN halaman_ttd DROP DEFAULT,
    ALTER COLUMN halaman_ttd TYPE smallint USING (CASE WHEN halaman_ttd THEN 1 ELSE 0 END);

COMMIT;