package view

import (
	"net/http"
)

func Render(rw http.ResponseWriter, code int, page string, data interface{}) {
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
