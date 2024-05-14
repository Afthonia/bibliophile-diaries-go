-- name: CreatePost :one
INSERT INTO posts (
    user_id,
    title,
    content
) VALUES (
    $1, $2, $3
) RETURNING id, 0::bigint as author_id, '' as author, COALESCE(title, '')::text as title, content, false as is_liked, 0 as like_count, created_at, '' as image_url;

-- name: GetPost :one
SELECT p.id, u.id as author_id, u.name as author, COALESCE(p.title, '')::text as title, p.content, p.created_at, '' as image_url  FROM posts p
INNER JOIN users u ON u.id = p.user_id
WHERE p.id = $1 LIMIT 1;

-- --name: GetUserPosts :many
-- SELECT p.id, u.id as author_id, u.name as author, COALESCE(p.title, '')::text as title, p.content, COALESCE(pl.is_liked, false)::bool as is_liked, COALESCE(pl2.like_count, 0)::int as like_count, p.created_at FROM posts p
-- INNER JOIN users u ON u.id = p.user_id
-- LEFT JOIN post_likes ON pl.post_id = p.id AND pl.user_id = u.id;

-- name: ListPosts :many
SELECT p.id, u.id as author_id, u.name as author, COALESCE(p.title, '')::text as title, p.content, COALESCE(pl.is_liked, false)::bool as is_liked, COALESCE(pl2.like_count, 0)::int as like_count, p.created_at, '' as image_url FROM posts p
INNER JOIN users u ON u.id = p.user_id
LEFT JOIN post_likes pl ON pl.post_id = p.id AND pl.user_id = $1
LEFT JOIN (SELECT post_id, COUNT(*) as like_count FROM post_likes WHERE is_liked = true GROUP BY post_id) pl2 ON pl2.post_id = p.id
ORDER BY p.created_at DESC;

-- name: GetLikedPosts :many
SELECT p.id, u.id as author_id, u.name as author, COALESCE(p.title, '')::text as title, p.content, COALESCE(pl.is_liked, false)::bool as is_liked, COALESCE(pl2.like_count, 0)::int as like_count, p.created_at, '' as image_url FROM posts p
INNER JOIN users u ON u.id = p.user_id
INNER JOIN post_likes pl ON pl.post_id = p.id AND pl.user_id = $1 and pl.is_liked = true
LEFT JOIN (SELECT post_id, COUNT(*) as like_count FROM post_likes WHERE is_liked = true GROUP BY post_id) pl2 ON pl2.post_id = p.id
ORDER BY p.created_at DESC;

-- -- name: GetPostLikes :many
-- SELECT p.id, COUNT(pl.*) FROM posts p
-- LEFT JOIN post_likes pl ON pl.post_id = p.id AND pl.is_liked = true


-- name: DeletePost :one
DELETE FROM posts
WHERE user_id = $1 AND id = $2
RETURNING id;

-- name: UpdatePost :one
UPDATE posts
  set title = $3,
  content = $4
WHERE user_id = $1 and id = $2
RETURNING id, 0::bigint as author_id, '' as author, COALESCE(title, '')::text as title, content, false as is_liked, 0 as like_count, created_at, '' as image_url;