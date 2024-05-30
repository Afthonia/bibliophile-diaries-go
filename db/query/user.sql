-- name: CreateUser :one
INSERT INTO users (
    permgroup_id,
    name,
    email,
    password_hash
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT id, name, email, created_at FROM users
ORDER BY id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateBio :one
UPDATE users
    set bio = $2
WHERE id = $1
RETURNING bio;

-- name: GetDashboard :one
SELECT 
    (SELECT COUNT(*) FROM posts WHERE user_id = $1) AS user_posts,
    (SELECT COUNT(*) FROM comments WHERE user_id = $1) AS user_comments;
