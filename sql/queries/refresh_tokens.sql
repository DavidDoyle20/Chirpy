-- name: AssignRefreshTokenToUser :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $2,
    NOW(),
    NOW(),
    $1,
    (NOW() + interval '60 day'),
    NULL
)
RETURNING *;

-- name: RevokeRefreshTokenFromUser :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM refresh_tokens
JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: CheckAndFetchRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1
    AND expires_at > NOW()
    AND revoked_at IS NULL;

