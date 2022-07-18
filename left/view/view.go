package view

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed template
var embedFS embed.FS
var templates *template.Template

func init() {
	templateFS, err := fs.Sub(embedFS, "template")
	if err != nil {
		panic(err)
	}
	templates = template.Must(template.New("").Funcs(helperMap).ParseFS(templateFS, "*.html", "**/*.html"))
}

func Render(rw http.ResponseWriter, code int, data interface{}, page string) {
	rw.WriteHeader(code)
	templates.ExecuteTemplate(rw, page, data)
}

func RenderError(rw http.ResponseWriter, code int, err error) {
	rw.WriteHeader(code)
	templates.ExecuteTemplate(rw, ErrorPage, struct {
		Code  int
		Error error
	}{code, err})
}
