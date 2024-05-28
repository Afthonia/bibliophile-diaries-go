package models

import (
	"errors"
	"net/http"
	"strconv"
)

type PostBind struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	BookTitle string `json:"book_title"`
	Vote      int    `json:"vote"`
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
	p.BookTitle = r.FormValue("book_title")
	vote, err := strconv.Atoi(r.FormValue("vote"))
	if err != nil {
		panic(err)
	}

	p.Vote = vote

	if len(p.Content) == 0 {
		return errors.New("your post needs to have a content")
	}

	return nil
}
