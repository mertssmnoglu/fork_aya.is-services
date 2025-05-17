-- name: GetStoryIdBySlug :one
SELECT id
FROM "story"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetStoryById :one
SELECT sqlc.embed(s), sqlc.embed(st)
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.id = sqlc.arg(id)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: GetStoryBySlug :one
SELECT sqlc.embed(s), sqlc.embed(st)
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.slug = sqlc.arg(slug)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: ListStories :many
SELECT sqlc.embed(s), sqlc.embed(st)
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR s.kind = sqlc.narg(filter_kind)::TEXT)
  AND (sqlc.narg(filter_author_profile_id)::CHAR(26) IS NULL OR s.author_profile_id = sqlc.narg(filter_author_profile_id)::CHAR(26))
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.deleted_at IS NULL
ORDER BY s.published_at DESC;
