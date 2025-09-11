BEGIN;

CREATE TABLE ref_kelas_jabatan (
    id SERIAL PRIMARY KEY,
    kelas_jabatan TEXT,
    tunjangan_kinerja BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMIT;