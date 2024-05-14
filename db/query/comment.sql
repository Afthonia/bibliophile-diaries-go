-- name: CreateComment :one
INSERT INTO comments (
    post_id,
    user_id,
    content
) VALUES (
    $1, $2, $3
) RETURNING id, post_id, 0::bigint AS commenter_id, '' AS commenter, content, created_at;

-- name: GetComment :one
SELECT c.id, u.id as commenter_id, u.name as commenter, c.content, c.created_at FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE c.id = $1 LIMIT 1;

-- name: GetCommentIDs :many
SELECT c.id AS comment_ids FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE post_id = $1
ORDER BY c.created_at DESC;

-- name: ListPostComments :many
SELECT c.id, c.post_id, u.id as commenter_id, u.name as commenter, c.content, c.created_at FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE post_id = $1
ORDER BY created_at;

-- name: ListUserComments :many
SELECT c.id, c.post_id, u.id as commenter_id, u.name as commenter, c.content, c.created_at  FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE user_id = $1
ORDER BY created_at;

-- name: DeleteComment :one
DELETE FROM comments
WHERE user_id = $1 AND id = $2
RETURNING id;

-- name: UpdateComment :one
UPDATE comments
  set content = $3
WHERE id = $1 and user_id = $2
RETURNING id, post_id, 0::bigint AS commenter_id, '' AS commenter, content, created_at;