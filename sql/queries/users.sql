-- name: CrateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;


-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;


-- name: GetUserById :one
SELECT id, created_at, updated_at, email, is_chirpy_red FROM users
WHERE id = $1;


-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM users
JOIN refresh_tokens
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1
  AND refresh_tokens.expires_at > NOW()
  AND refresh_tokens.revoked_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET email = $1 , password = $2 , updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red; 


-- name: UpdateUserChirpy :one
UPDATE users
SET is_chirpy_red = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;