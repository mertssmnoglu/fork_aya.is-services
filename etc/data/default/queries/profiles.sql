-- name: GetProfileIdBySlug :one
SELECT id
FROM "profile"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetProfileIdByCustomDomain :one
SELECT id
FROM "profile"
WHERE custom_domain = sqlc.arg(custom_domain)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetProfileById :one
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.id = sqlc.arg(id)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: ListProfiles :many
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE (sqlc.narg(filter_kind)::TEXT IS NULL OR p.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND p.deleted_at IS NULL;

-- name: CreateProfile :one
INSERT INTO "profile" (id, slug)
VALUES (sqlc.arg(id), sqlc.arg(slug)) RETURNING *;

-- name: UpdateProfile :execrows
UPDATE "profile"
SET slug = sqlc.arg(slug)
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: RemoveProfile :execrows
UPDATE "profile"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetProfileLinksForKind :many
SELECT pl.*
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id
  AND p.deleted_at IS NULL
WHERE pl.kind = sqlc.arg(kind)
  AND pl.deleted_at IS NULL
ORDER BY pl."order";

-- name: GetProfilePagesByProfileId :many
SELECT pp.*, ppt.*
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
  AND ppt.locale_code = sqlc.arg(locale_code)
WHERE pp.profile_id = sqlc.arg(profile_id)
  AND pp.deleted_at IS NULL
ORDER BY pp."order";

-- name: GetProfilePageByProfileIdAndSlug :one
SELECT pp.*, ppt.*
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
  AND ppt.locale_code = sqlc.arg(locale_code)
WHERE pp.profile_id = sqlc.arg(profile_id) AND pp.slug = sqlc.arg(page_slug) AND pp.deleted_at IS NULL
ORDER BY pp."order";

-- name: GetProfileLinksByProfileId :many
SELECT *
FROM "profile_link"
WHERE profile_id = sqlc.arg(profile_id)
  AND is_hidden = FALSE
  AND deleted_at IS NULL
ORDER BY "order";

-- name: ListProfileMembershipsByProfileIdAndKind :many
SELECT
  sqlc.embed(pm),
  sqlc.embed(pp),
  sqlc.embed(ppt)
FROM
	"profile_membership" pm
  INNER JOIN "profile" pp ON pp.id = pm.profile_id AND pp.kind = ANY(string_to_array(sqlc.arg(kind)::TEXT, ',')) AND pp.deleted_at IS NULL
  INNER JOIN "profile_tx" ppt ON ppt.profile_id = pp.id
	  AND ppt.locale_code = sqlc.arg(locale_code)
  INNER JOIN "user" u ON u.id = pm.user_id AND u.deleted_at IS NULL
  INNER JOIN "profile" pc ON pc.id = u.individual_profile_id AND pc.deleted_at IS NULL
WHERE pc.id = sqlc.arg(profile_id)
  AND pm.deleted_at IS NULL;
