-- name: GetProfileById :one
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON p.id = pt.profile_id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.id = sqlc.arg(id)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: GetProfileBySlug :one
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON p.id = pt.profile_id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.slug = sqlc.arg(slug)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: GetProfileByCustomDomain :one
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON p.id = pt.profile_id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.custom_domain = sqlc.arg(domain)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: ListProfiles :many
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON p.id = pt.profile_id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.deleted_at IS NULL;

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
