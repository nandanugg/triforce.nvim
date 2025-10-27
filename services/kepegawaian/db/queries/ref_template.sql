-- name: ListTemplates :many
SELECT id, nama, filename, created_at, updated_at
FROM ref_template
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountTemplates :one
SELECT COUNT(1)
FROM ref_template
WHERE deleted_at IS NULL;

-- name: GetTemplate :one
SELECT id, nama, filename, created_at, updated_at
FROM ref_template
WHERE id = @id::integer AND deleted_at IS NULL;

-- name: GetTemplateBerkas :one
SELECT file_base64
FROM ref_template
WHERE id = @id::integer AND deleted_at IS NULL;

-- name: CreateTemplate :one
INSERT INTO ref_template (nama, filename, file_base64)
VALUES (@name::text, @filename::text, @file_base64::text)
RETURNING id, nama, filename, created_at, updated_at;

-- name: UpdateTemplate :one
UPDATE ref_template
SET
	nama = @name::text,
	filename = @filename::text,
	file_base64 = @file_base64::text,
	updated_at = NOW()
WHERE id = @id::integer AND deleted_at IS NULL
RETURNING id, nama, filename, created_at, updated_at;

-- name: DeleteTemplate :execrows
UPDATE ref_template
SET deleted_at = NOW()
WHERE id = @id::integer AND deleted_at IS NULL;

