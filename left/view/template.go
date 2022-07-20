package view

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

func Render(rw http.ResponseWriter, code int, page string, data interface{}) {
	rw.WriteHeader(code)
	getTemplate(page).Execute(rw, data)
}

func parseTemplate(name string, templateFS fs.FS) *template.Template {
	return template.Must(template.New(name).Funcs(helperMap).ParseFS(templateFS, name, "**/*.html"))
}

const (
	IndexPage    string = "index.html"
	EnvelopePage string = "envelope.html"
)

type IndexData struct {
	Envelopes []envelope.Envelope
	Page      paginate.Page
}

type EnvelopeData struct {
	Envelope *envelope.Envelope
	Tab      string
}
