-- name: UpdateRiwayatUjiKompetensiNamaNipByNIP :exec
UPDATE riwayat_uji_kompetensi
SET
    nip_baru = @nip_baru::varchar,
    updated_at = now()
WHERE nip_baru = @nip::varchar AND deleted_at IS NULL
AND (
    (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM nip_baru)
);