-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, email, password_hash)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: UpdateAPIKey :one
UPDATE users 
SET api_key = $2
WHERE id = $1 
RETURNING api_key;

-- name: GetUserStats :one
SELECT 
    u.username,
    COUNT(tc.id) as completed_tasks
FROM users u
LEFT JOIN task_completions tc ON u.id = tc.user_id
WHERE u.id = $1
GROUP BY u.username;
