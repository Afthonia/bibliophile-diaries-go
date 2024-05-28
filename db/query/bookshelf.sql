-- name: ToggleBook :one
INSERT INTO bookshelf (
    book_id,
    user_id
) VALUES (
    $1, $2
) ON CONFLICT(book_id, user_id) DO UPDATE SET
in_bookshelf = NOT bookshelf.in_bookshelf
RETURNING in_bookshelf;

