-- name: TogglePostLike :one
INSERT INTO post_likes (
    post_id,
    user_id
) VALUES (
    $1, $2
) ON CONFLICT(post_id,user_id) DO UPDATE SET
is_liked = NOT post_likes.is_liked
RETURNING is_liked;
