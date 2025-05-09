-- name: GetStoryById :one
SELECT *
FROM "story"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetStoryBySlug :one
SELECT *
FROM "story"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListStories :many
SELECT *
FROM "story"
WHERE deleted_at IS NULL;
