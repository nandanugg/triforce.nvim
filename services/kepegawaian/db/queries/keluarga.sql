-- name: ListOrangTuaByNip :many
SELECT
    ot.id,
    ot.hubungan,
    ot.nama,
    ot.tanggal_meninggal,
    ot.akte_meninggal,
    ot.no_dokumen AS nik,
    ot.agama_id,
    ra.nama AS agama
FROM orang_tua ot
JOIN pegawai pg ON ot.pns_id = pg.pns_id AND pg.deleted_at is null
LEFT JOIN ref_agama ra ON ot.agama_id = ra.id AND ra.deleted_at is null
WHERE ot.deleted_at IS NULL
AND pg.nip_baru = $1;

-- name: ListPasanganByNip :many
SELECT
    p.id,
    p.pns,
    p.nama,
    p.nik,
    p.karsus,
    p.status as status_pernikahan_id,
    rjk.nama as status_pernikahan,
    p.agama_id,
    ra.nama AS agama,
    p.tanggal_lahir,
    p.tanggal_menikah,
    p.tanggal_cerai,
    p.tanggal_meninggal,
    p.akte_nikah,
    p.akte_cerai,
    p.akte_meninggal
FROM pasangan p
JOIN pegawai pg ON p.pns_id = pg.pns_id AND pg.deleted_at is null
LEFT JOIN ref_agama ra ON p.agama_id = ra.id AND ra.deleted_at is null
LEFT JOIN ref_jenis_kawin rjk ON p.status = rjk.id AND rjk.deleted_at is null
WHERE p.deleted_at IS NULL
AND pg.nip_baru = $1;

-- name: ListAnakByNip :many
SELECT
    a.id,
    a.pasangan_id,
    a.nama,
    a.jenis_kelamin,
    a.tanggal_lahir,
    a.status_anak,
    pas.nama AS nama_orang_tua,
    a.nik,
    a.agama_id,
    ra.nama as agama,
    a.jenis_kawin_id,
    rjk.nama as status_pernikahan,
    a.status_sekolah,
    a.anak_ke
FROM anak a
JOIN pegawai pg ON a.pns_id = pg.pns_id AND pg.deleted_at is null
LEFT JOIN pasangan pas ON a.pasangan_id = pas.id AND pas.deleted_at is null
LEFT JOIN ref_agama ra ON a.agama_id = ra.id AND ra.deleted_at is null
LEFT JOIN ref_jenis_kawin rjk ON a.jenis_kawin_id = rjk.id AND rjk.deleted_at is null
WHERE a.deleted_at IS NULL
AND pg.nip_baru = $1
ORDER BY a.anak_ke NULLS LAST, a.tanggal_lahir NULLS LAST;

-- name: CreateOrangTua :one
insert into orang_tua
    (nama, jenis_dokumen, no_dokumen, hubungan, agama_id, tanggal_meninggal, akte_meninggal, pns_id, nip) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning id;

-- name: UpdateOrangTua :execrows
update orang_tua ot
set
    nama = $1,
    jenis_dokumen = $2,
    no_dokumen = $3,
    hubungan = $4,
    agama_id = $5,
    tanggal_meninggal = $6,
    akte_meninggal = $7,
    updated_at = now()
from pegawai p
where ot.pns_id = p.pns_id
    and ot.id = @id and ot.deleted_at is null
    and p.nip_baru = @nip::varchar and p.deleted_at is null;

-- name: DeleteOrangTua :execrows
update orang_tua ot
set deleted_at = now()
from pegawai p
where ot.pns_id = p.pns_id
    and ot.id = @id and ot.deleted_at is null
    and p.nip_baru = @nip::varchar and p.deleted_at is null;

-- name: CreatePasangan :one
insert into pasangan
    (nama, nik, pns, tanggal_lahir, karsus, agama_id, status, hubungan, tanggal_menikah, akte_nikah, tanggal_meninggal, akte_meninggal, tanggal_cerai, akte_cerai, pns_id, nip) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
returning id;

-- name: IsPasanganExistsByIDAndNIP :one
select exists (
  select 1 from pasangan ps
  join pegawai p on p.pns_id = ps.pns_id and p.deleted_at is null
  where ps.id = @id and p.nip_baru = @nip::varchar and ps.deleted_at is null
);

-- name: UpdatePasangan :execrows
update pasangan ps
set
    nama = $1,
    nik = $2,
    pns = $3,
    tanggal_lahir = $4,
    karsus = $5,
    agama_id = $6,
    status = $7,
    hubungan = $8,
    tanggal_menikah = $9,
    akte_nikah = $10,
    tanggal_meninggal = $11,
    akte_meninggal = $12,
    tanggal_cerai = $13,
    akte_cerai = $14,
    updated_at = now()
from pegawai p
where ps.pns_id = p.pns_id
    and ps.id = @id and ps.deleted_at is null
    and p.nip_baru = @nip::varchar and p.deleted_at is null;

-- name: DeletePasangan :execrows
update pasangan ps
set deleted_at = now()
from pegawai p
where ps.pns_id = p.pns_id
    and ps.id = @id and ps.deleted_at is null
    and p.nip_baru = @nip::varchar and p.deleted_at is null;

-- name: CreateAnak :one
insert into anak
    (pasangan_id, nama, nik, jenis_kelamin, tanggal_lahir, agama_id, jenis_kawin_id, status_anak, status_sekolah, anak_ke, pns_id, nip) values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
returning id;

-- name: UpdateAnak :execrows
update anak a
set
    pasangan_id = $1,
    nama = $2,
    nik = $3,
    jenis_kelamin = $4,
    tanggal_lahir = $5,
    agama_id = $6,
    jenis_kawin_id = $7,
    status_anak = $8,
    status_sekolah = $9,
    anak_ke = $10,
    updated_at = now()
from pegawai p
where a.pns_id = p.pns_id
    and a.id = @id and a.deleted_at is null
    and p.nip_baru = @nip::varchar and p.deleted_at is null;

-- name: DeleteAnak :execrows
update anak a
set deleted_at = now()
from pegawai p
where a.pns_id = p.pns_id
    and a.id = @id and a.deleted_at is null
    and p.nip_baru = @nip::varchar and p.deleted_at is null;

-- name: UpdateAnakNIPByNIP :exec
UPDATE anak
SET 
    nip = @nip_baru::varchar,
    updated_at = now()
WHERE nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip)
);

-- name: UpdateOrangTuaNIPByNIP :exec
UPDATE orang_tua
SET 
    nip = @nip_baru::varchar,
    updated_at = now()
WHERE nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip)
);

-- name: UpdatePasanganNIPByNIP :exec
UPDATE pasangan
SET 
    nip = @nip_baru::varchar,
    updated_at = now()
WHERE nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip)
);
