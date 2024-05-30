-- name: ToggleBook :one
INSERT INTO bookshelf (
    book_id,
    user_id
) VALUES (
    $1, $2
) ON CONFLICT(book_id, user_id) DO UPDATE SET
in_bookshelf = NOT bookshelf.in_bookshelf
RETURNING in_bookshelf;

-- name: ListBookshelf :many
SELECT user_id, book_id, in_bookshelf, created_at FROM bookshelf
WHERE user_id = $1
ORDER BY created_at DESC;