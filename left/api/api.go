package api

import (
	"net/http"
)

type (
	Response struct {
		Code  int
		Data  interface{}
		Error error
	}

	Handler func(w http.ResponseWriter, r *http.Request) Response

	Renderer interface {
		Render(rw http.ResponseWriter, r Response)
	}
)
