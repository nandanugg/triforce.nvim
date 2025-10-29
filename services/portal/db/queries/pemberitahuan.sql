-- name: CreatePemberitahuan :one
INSERT INTO pemberitahuan (
    judul_berita,
    deskripsi_berita,
    pinned,
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
    pinned,
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

-- name: GetOverlappingPinnedPemberitahuan :one
SELECT
    id,
    judul_berita
FROM pemberitahuan
WHERE
    deleted_at IS NULL
    AND pinned = TRUE
    AND diterbitkan_pada <= @ditarik_pada::timestamptz
    AND ditarik_pada >= @diterbitkan_pada::timestamptz
    AND (
        @id::bigint = 0 OR id <> @id::bigint
    )
LIMIT 1;

-- name: ListPemberitahuan :many
SELECT
    id,
    judul_berita,
    deskripsi_berita,
    pinned,
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
    END AS status
FROM pemberitahuan
WHERE
    deleted_at IS NULL
    AND (
        @status = 'ALL'
        OR (@status = 'WAITING' AND NOW() < diterbitkan_pada)
        OR (@status = 'ACTIVE' AND NOW() >= diterbitkan_pada AND NOW() < ditarik_pada)
        OR (@status = 'OVER' AND NOW() >= ditarik_pada)
    )
ORDER BY
    (pinned = TRUE AND NOW() >= diterbitkan_pada AND NOW() < ditarik_pada) DESC,
    diterbitkan_pada DESC
LIMIT $1 OFFSET $2;

-- name: CountPemberitahuan :one
SELECT COUNT(1) AS total
FROM pemberitahuan
WHERE
    deleted_at IS NULL
    AND (
        @status = 'ALL'
        OR (@status = 'WAITING' AND NOW() < diterbitkan_pada)
        OR (@status = 'ACTIVE' AND NOW() >= diterbitkan_pada AND NOW() < ditarik_pada)
        OR (@status = 'OVER' AND NOW() >= ditarik_pada)
    );

-- name: UpdatePemberitahuan :one
UPDATE pemberitahuan
SET
    judul_berita = $2,
    deskripsi_berita = $3,
    pinned = $4,
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
    pinned,
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

