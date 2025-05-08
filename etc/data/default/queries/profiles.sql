-- name: GetProfileById :one
SELECT *
FROM "profile"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetProfileBySlug :one
SELECT *
FROM "profile"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListProfiles :many
SELECT *
FROM "profile"
WHERE deleted_at IS NULL;

-- name: CreateProfile :one
INSERT INTO "profile" (id, slug)
VALUES (sqlc.arg(id), sqlc.arg(slug)) RETURNING *;

-- name: UpdateProfile :execrows
UPDATE "profile"
SET slug = sqlc.arg(slug)
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DeleteProfile :execrows
UPDATE "profile"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetProfileLinksForKind :many
SELECT pl.*
FROM "profile_link" pl
  INNER JOIN "profile" p ON pl.profile_id = p.id
  AND p.deleted_at IS NULL
WHERE pl.kind = sqlc.arg(kind)
  AND pl.deleted_at IS NULL;
