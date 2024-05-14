package models

import (
	"errors"
	"net/http"
)

type CommentBind struct {
	PostID  int64  `json:"post_id"`
	Content string `json:"content"`
}

func (c *CommentBind) Bind(r *http.Request) error {
	if len(c.Content) == 0 {
		return errors.New("your comment needs to have a content")
	}

	return nil
}
