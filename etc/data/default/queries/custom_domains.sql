-- name: GetCustomDomainByDomain :one
SELECT cd.*,
  p.slug AS "profile_slug"
FROM "custom_domain" cd
  INNER JOIN "profile" p ON cd.profile_id = p.id
  AND p.deleted_at IS NULL
WHERE cd.domain = sqlc.arg(domain)
  AND cd.deleted_at IS NULL
LIMIT 1;
