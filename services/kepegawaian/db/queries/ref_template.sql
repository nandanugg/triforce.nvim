-- name: ListTemplates :many
SELECT id, nama, created_at, updated_at
FROM ref_template
WHERE deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: CountTemplates :one
SELECT COUNT(1)
FROM ref_template
WHERE deleted_at IS NULL;

-- name: GetTemplate :one
SELECT id, nama, created_at, updated_at
FROM ref_template
WHERE id = @id::integer AND deleted_at IS NULL;

-- name: GetTemplateBerkas :one
SELECT file_base64
FROM ref_template
WHERE id = @id::integer AND deleted_at IS NULL;

-- name: CreateTemplate :one
INSERT INTO ref_template (nama, file_base64)
VALUES (@name::text, @file_base64::text)
RETURNING id, nama, created_at, updated_at;

-- name: UpdateTemplate :one
UPDATE ref_template
SET nama = @name::text,
    file_base64 = @file_base64::text,
    updated_at = NOW()
WHERE id = @id::integer AND deleted_at IS NULL
RETURNING id, nama, created_at, updated_at;

-- name: DeleteTemplate :execrows
UPDATE ref_template
SET deleted_at = NOW()
WHERE id = @id::integer AND deleted_at IS NULL;

