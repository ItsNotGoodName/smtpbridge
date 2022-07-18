package view

import (
	"net/http"
)

func Render(rw http.ResponseWriter, code int, data interface{}, page string) {
	rw.WriteHeader(code)
	getTemplate(page).Execute(rw, data)
}

func RenderError(rw http.ResponseWriter, code int, err error) {
	rw.WriteHeader(code)
	getTemplate(ErrorPage).Execute(rw, struct {
		Code  int
		Error error
	}{code, err})
}
