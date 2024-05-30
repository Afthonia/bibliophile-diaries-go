package db

import (
	"net/http"
)

func (l *ListPostsRow) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (l *ListPostCommentsRow) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (l *ListUserCommentsRow) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (l *GetPostRow) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (l *ListBookshelfRow) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
