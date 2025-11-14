-- name: CreatePemberitahuan :one
INSERT INTO pemberitahuan (
    judul_berita,
    deskripsi_berita,
    pinned_at,
    diterbitkan_pada,
    ditarik_pada,
    updated_by,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, now()
)
RETURNING
    id,
    judul_berita,
    deskripsi_berita,
    pinned_at,
    diterbitkan_pada,
    ditarik_pada,
    updated_by,
    updated_at,
    CASE
        WHEN NOW() < diterbitkan_pada THEN 'WAITING'
        WHEN NOW() >= diterbitkan_pada AND NOW() < ditarik_pada THEN 'ACTIVE'
        WHEN NOW() >= ditarik_pada THEN 'OVER'
        ELSE 'UNKNOWN'
    END AS status;

-- name: ListActivePemberitahuan :many
WITH active_pemberitahuan AS (
    SELECT
        id,
        judul_berita,
        deskripsi_berita,
        pinned_at,
        diterbitkan_pada,
        ditarik_pada,
        updated_by,
        updated_at,
        ROW_NUMBER() OVER (
            ORDER BY pinned_at DESC NULLS LAST
        ) AS pinned_rank
    FROM pemberitahuan
    WHERE
        deleted_at IS NULL
        AND aktif_range @> now()
)
SELECT
    id,
    judul_berita,
    deskripsi_berita,
    pinned_at,
    diterbitkan_pada,
    ditarik_pada,
    updated_by,
    updated_at,
    (pinned_at IS NOT NULL AND pinned_rank = 1) AS is_current_period_pinned
FROM active_pemberitahuan
ORDER BY
    is_current_period_pinned DESC,
    diterbitkan_pada DESC
LIMIT $1 OFFSET $2;


-- name: ListPemberitahuan :many
SELECT
    id,
    judul_berita,
    deskripsi_berita,
    pinned_at,
    diterbitkan_pada,
    ditarik_pada,
    updated_by,
    updated_at,
    aktif_range
FROM pemberitahuan
WHERE
    deleted_at IS NULL
    AND (
        @judul_berita = '' OR judul_berita ILIKE CONCAT('%', @judul_berita, '%')
    )
ORDER BY
    CASE
        WHEN @sort_by = 'pinned_asc' THEN pinned_at
    END ASC NULLS LAST,
    CASE
        WHEN @sort_by = 'pinned_desc' THEN pinned_at
    END DESC NULLS LAST,
    diterbitkan_pada DESC
LIMIT $1 OFFSET $2;

-- name: CountPemberitahuan :one
SELECT COUNT(1) AS total
FROM pemberitahuan
WHERE
    deleted_at IS NULL
    AND (
        @status = 'ALL'
        OR (@status = 'ACTIVE' AND aktif_range @> now())
    )
    AND (
        @judul_berita = '' OR judul_berita ILIKE CONCAT('%', @judul_berita, '%')
    );

-- name: UpdatePemberitahuan :one
UPDATE pemberitahuan
SET
    judul_berita = $2,
    deskripsi_berita = $3,
    pinned_at = $4,
    diterbitkan_pada = $5,
    ditarik_pada = $6,
    updated_by = $7,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING
    id,
    judul_berita,
    deskripsi_berita,
    pinned_at,
    diterbitkan_pada,
    ditarik_pada,
    updated_by,
    updated_at,
    deleted_at,
    CASE
        WHEN NOW() < diterbitkan_pada THEN 'WAITING'
        WHEN NOW() >= diterbitkan_pada AND NOW() < ditarik_pada THEN 'ACTIVE'
        WHEN NOW() >= ditarik_pada THEN 'OVER'
        ELSE 'UNKNOWN'
    END AS status;

-- name: DeletePemberitahuan :execrows
UPDATE pemberitahuan
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

