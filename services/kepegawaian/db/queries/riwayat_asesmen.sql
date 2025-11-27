-- name: UpdateRiwayatAsesmenNamaNipByPNSID :exec
UPDATE riwayat_asesmen
SET     
    pns_nip = @nip_baru::varchar,
    nama_lengkap = @nama::varchar,
    updated_at = now()
WHERE 
    pns_id = @pns_id::varchar AND deleted_at IS NULL
    AND (
        (@nip_baru::varchar IS NOT NULL AND @nip_baru::varchar IS DISTINCT FROM pns_nip)
        OR (@nama::varchar IS NOT NULL AND @nama::varchar IS DISTINCT FROM nama_lengkap)
    );