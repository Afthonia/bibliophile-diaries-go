// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: comment.sql

package db

import (
	"context"
	"time"
)

const createComment = `-- name: CreateComment :one
INSERT INTO comments (
    post_id,
    user_id,
    content
) VALUES (
    $1, $2, $3
) RETURNING id, post_id, 0::bigint AS commenter_id, '' AS commenter, content, created_at
`

type CreateCommentParams struct {
	PostID  int64  `json:"post_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}

type CreateCommentRow struct {
	ID          int64     `json:"id"`
	PostID      int64     `json:"post_id"`
	CommenterID int64     `json:"commenter_id"`
	Commenter   string    `json:"commenter"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (CreateCommentRow, error) {
	row := q.db.QueryRowContext(ctx, createComment, arg.PostID, arg.UserID, arg.Content)
	var i CreateCommentRow
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.CommenterID,
		&i.Commenter,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :one
DELETE FROM comments
WHERE user_id = $1 AND id = $2
RETURNING id
`

type DeleteCommentParams struct {
	UserID int64 `json:"user_id"`
	ID     int64 `json:"id"`
}

func (q *Queries) DeleteComment(ctx context.Context, arg DeleteCommentParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteComment, arg.UserID, arg.ID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getComment = `-- name: GetComment :one
SELECT c.id, u.id as commenter_id, u.name as commenter, c.content, c.created_at FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE c.id = $1 LIMIT 1
`

type GetCommentRow struct {
	ID          int64     `json:"id"`
	CommenterID int64     `json:"commenter_id"`
	Commenter   string    `json:"commenter"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) GetComment(ctx context.Context, id int64) (GetCommentRow, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var i GetCommentRow
	err := row.Scan(
		&i.ID,
		&i.CommenterID,
		&i.Commenter,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const getCommentIDs = `-- name: GetCommentIDs :many
SELECT c.id AS comment_ids FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE post_id = $1
ORDER BY c.created_at DESC
`

func (q *Queries) GetCommentIDs(ctx context.Context, postID int64) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getCommentIDs, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var comment_ids int64
		if err := rows.Scan(&comment_ids); err != nil {
			return nil, err
		}
		items = append(items, comment_ids)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPostComments = `-- name: ListPostComments :many
SELECT c.id, c.post_id, u.id as commenter_id, u.name as commenter, c.content, c.created_at FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE post_id = $1
ORDER BY created_at
`

type ListPostCommentsRow struct {
	ID          int64     `json:"id"`
	PostID      int64     `json:"post_id"`
	CommenterID int64     `json:"commenter_id"`
	Commenter   string    `json:"commenter"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) ListPostComments(ctx context.Context, postID int64) ([]ListPostCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, listPostComments, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPostCommentsRow
	for rows.Next() {
		var i ListPostCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.CommenterID,
			&i.Commenter,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUserComments = `-- name: ListUserComments :many
SELECT c.id, c.post_id, u.id as commenter_id, u.name as commenter, c.content, c.created_at  FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE user_id = $1
ORDER BY created_at
`

type ListUserCommentsRow struct {
	ID          int64     `json:"id"`
	PostID      int64     `json:"post_id"`
	CommenterID int64     `json:"commenter_id"`
	Commenter   string    `json:"commenter"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) ListUserComments(ctx context.Context, userID int64) ([]ListUserCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserComments, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserCommentsRow
	for rows.Next() {
		var i ListUserCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.CommenterID,
			&i.Commenter,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateComment = `-- name: UpdateComment :one
UPDATE comments
  set content = $3
WHERE id = $1 and user_id = $2
RETURNING id, post_id, 0::bigint AS commenter_id, '' AS commenter, content, created_at
`

type UpdateCommentParams struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}

type UpdateCommentRow struct {
	ID          int64     `json:"id"`
	PostID      int64     `json:"post_id"`
	CommenterID int64     `json:"commenter_id"`
	Commenter   string    `json:"commenter"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) UpdateComment(ctx context.Context, arg UpdateCommentParams) (UpdateCommentRow, error) {
	row := q.db.QueryRowContext(ctx, updateComment, arg.ID, arg.UserID, arg.Content)
	var i UpdateCommentRow
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.CommenterID,
		&i.Commenter,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}