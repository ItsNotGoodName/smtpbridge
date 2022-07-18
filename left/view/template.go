package view

import (
	"html/template"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

func parseTemplate(name string, templateFS fs.FS) *template.Template {
	return template.Must(template.New(name).Funcs(helperMap).ParseFS(templateFS, name, "**/*.html"))
}

const (
	ErrorPage    string = "error.html"
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
