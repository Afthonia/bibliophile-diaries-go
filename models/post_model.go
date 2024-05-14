package models

import (
	"errors"
	"net/http"
)

type PostBind struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (p *PostBind) Bind(r *http.Request) error {

	if len(p.Content) == 0 {
		return errors.New("your post needs to have a content")
	}

	return nil
}

func (p *PostBind) FormBind(r *http.Request) error {
	p.Title = r.FormValue("title")
	p.Content = r.FormValue("content")

	if len(p.Content) == 0 {
		return errors.New("your post needs to have a content")
	}

	return nil
}
