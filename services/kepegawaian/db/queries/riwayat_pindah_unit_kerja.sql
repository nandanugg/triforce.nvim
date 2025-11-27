-- name: UpdateRiwayatPindahUnitKerjaNamaNipByNIP :exec
UPDATE riwayat_pindah_unit_kerja
SET     
    pns_nip = @nip_baru::varchar,
    pns_nama = @nama::varchar,
    updated_at = now()
WHERE pns_nip = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM pns_nip)
    OR (@nama::varchar IS NOT NULL AND @nama::varchar IS DISTINCT FROM pns_nama)
);