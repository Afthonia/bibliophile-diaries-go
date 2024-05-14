package status

import (
	"database/sql"
	"net/http"

	"log"

	"github.com/go-chi/render"
)

type StatusResponse struct {
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *StatusResponse) Render(w http.ResponseWriter, r *http.Request) error {
	log.Println(e.ErrorText)
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrBadRequest(err error) render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 400,
		StatusText:     "Bad request.",
		ErrorText:      err.Error(),
	}
}

func ErrInternal(err error) render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

func ErrNotFoundOrInternal(err error) render.Renderer {
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return &StatusResponse{
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrConflict(err string) render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 409,
		StatusText:     "Conflict.",
		ErrorText:      err,
	}
}

func ErrUnauthorized(err string) render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 401,
		StatusText:     "Unauthorized.",
		ErrorText:      err,
	}
}

func DelSuccess() render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 200,
		StatusText:     "Successfuly deleted.",
	}
}

func Success() render.Renderer {
	return &StatusResponse{
		HTTPStatusCode: 200,
		StatusText:     "Success !",
	}
}

var ErrNotFound = &StatusResponse{HTTPStatusCode: 404, StatusText: "Not found."}

type SuccesWithID struct {
	StatusText string `json:"status"`
	ID         int64  `json:"id"`
}

func SuccessID(id int64, text string) render.Renderer {
	return &SuccesWithID{
		StatusText: text,
		ID:         id,
	}
}

func (d *SuccesWithID) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, 200)
	return nil
}

type SuccesWithCount struct {
	StatusText string `json:"status"`
	Count      int64  `json:"count"`
}

func SuccessCount(count int64, text string) render.Renderer {
	return &SuccesWithCount{
		StatusText: text,
		Count:      count,
	}
}

func (d *SuccesWithCount) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, 200)
	return nil
}

type IDWithApproved struct {
	ID       int64 `json:"id"`
	Approved bool  `json:"approved"`
}

func (d *IDWithApproved) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type IDWithRequestStatus struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func (d *IDWithRequestStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type OnlyStatus struct {
	Status string `json:"status"`
}

func (d *OnlyStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
