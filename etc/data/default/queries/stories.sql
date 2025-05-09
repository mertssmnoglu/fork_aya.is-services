-- name: GetStoryById :one
SELECT sqlc.embed(s), sqlc.embed(st)
FROM "story" s
  INNER JOIN "story_tx" st ON s.id = st.story_id
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.id = sqlc.arg(id)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: GetStoryBySlug :one
SELECT sqlc.embed(s), sqlc.embed(st)
FROM "story" s
  INNER JOIN "story_tx" st ON s.id = st.story_id
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.slug = sqlc.arg(slug)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: ListStories :many
SELECT sqlc.embed(s), sqlc.embed(st)
FROM "story" s
  INNER JOIN "story_tx" st ON s.id = st.story_id
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.deleted_at IS NULL;
