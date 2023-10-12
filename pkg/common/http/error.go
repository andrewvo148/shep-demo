package http

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HttpStatusCode int   `json:"-"` // http response status code

	AppCode   int64  `json:"app_code,omitempty"`   // application-specific error code
	ErrorText string `json:"error_text,omitempty"` // application-level error message
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HttpStatusCode)
	return nil
}

func ErrInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HttpStatusCode: http.StatusInternalServerError,
		ErrorText:      err.Error(),
	}
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HttpStatusCode: http.StatusBadRequest,
		ErrorText:      err.Error(),
	}
}
