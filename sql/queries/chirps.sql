-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ResetChirps :exec
DELETE FROM chirps;

-- name: GetChirp :one
SELECT * FROM chirps
    WHERE id = $1;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY chirps.created_at
;

-- name: GetChirpAuthor :one
SELECT users.* 
FROM chirps
JOIN users ON chirps.user_id = users.id
WHERE chirps.id = $1
;

-- name: RemoveChirp :exec
DELETE FROM chirps
WHERE id = $1
;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE chirps.user_id = $1
ORDER BY chirps.created_at
;