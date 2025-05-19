-- name: CreateSession :one
INSERT INTO session (
    id,
    status,
    oauth_request_state,
    oauth_request_code_verifier,
    oauth_redirect_uri,
    logged_in_user_id,
    logged_in_at,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, status, oauth_request_state, oauth_request_code_verifier, oauth_redirect_uri, logged_in_user_id, logged_in_at, expires_at, created_at, updated_at;

-- name: UpdateSessionLoggedInAt :exec
UPDATE session SET logged_in_at = $2, updated_at = NOW() WHERE id = $1;

-- name: GetSessionByID :one
SELECT id, status, oauth_request_state, oauth_request_code_verifier, oauth_redirect_uri, logged_in_user_id, logged_in_at, expires_at, created_at, updated_at
FROM session WHERE id = $1;
