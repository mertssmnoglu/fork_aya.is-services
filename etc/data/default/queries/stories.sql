-- name: GetStoryIDBySlug :one
SELECT id
FROM "story"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetStoryByID :one
SELECT sqlc.embed(s), sqlc.embed(st), sqlc.embed(p), sqlc.embed(pt)
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile" p ON p.id = s.author_profile_id AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = sqlc.arg(locale_code)
WHERE s.id = sqlc.arg(id)
  AND s.deleted_at IS NULL
LIMIT 1;

-- -- name: ListStories :many
-- SELECT sqlc.embed(s), sqlc.embed(st), sqlc.embed(p), sqlc.embed(pt)
-- FROM "story" s
--   INNER JOIN "story_tx" st ON st.story_id = s.id
--   AND (sqlc.narg(filter_kind)::TEXT IS NULL OR s.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
--   AND (sqlc.narg(filter_author_profile_id)::CHAR(26) IS NULL OR s.author_profile_id = sqlc.narg(filter_author_profile_id)::CHAR(26))
--   AND st.locale_code = sqlc.arg(locale_code)
--   LEFT JOIN "profile" p ON p.id = s.author_profile_id AND p.deleted_at IS NULL
--   INNER JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = sqlc.arg(locale_code)
-- WHERE s.deleted_at IS NULL
-- ORDER BY s.published_at DESC;

-- name: ListStories :many
SELECT
  sqlc.embed(s),
  sqlc.embed(st),
  sqlc.embed(p),
  sqlc.embed(pt)
FROM "story_publication" sp
  INNER JOIN "story" s ON s.id = sp.story_id
  AND s.deleted_at IS NULL
  AND s.published_at IS NOT NULL
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR s.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND (sqlc.narg(filter_author_profile_id)::CHAR(26) IS NULL OR s.author_profile_id = sqlc.narg(filter_author_profile_id)::CHAR(26))
  AND st.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile" p ON p.id = s.author_profile_id
  AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE
  (sqlc.narg(filter_publication_profile_id)::CHAR(26) IS NULL OR sp.profile_id = sqlc.narg(filter_publication_profile_id)::CHAR(26))
  AND s.deleted_at IS NULL
ORDER BY s.published_at DESC;
