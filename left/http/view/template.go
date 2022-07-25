package view

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
	"github.com/ItsNotGoodName/smtpbridge/core/version"
)

func Render(w http.ResponseWriter, code int, page string, data interface{}) {
	w.WriteHeader(code)
	getTemplate(page).Execute(w, data)
}

func parseTemplate(name string, templateFS fs.FS) *template.Template {
	return template.Must(template.New(name).Funcs(helperMap).ParseFS(templateFS, name, "**/*.html"))
}

const (
	IndexPage       string = "index.html"
	EnvelopePage    string = "envelope.html"
	EnvelopesPage   string = "envelopes.html"
	AttachmentsPage string = "attachments.html"
	EndpointsPage   string = "endpoints.html"
)

type IndexData struct {
	EnvelopesCount   int
	AttachmentsCount int
	Envelopes        []envelope.Envelope
	Version          version.Version
}

type EnvelopesData struct {
	Envelopes []envelope.Envelope
	Page      paginate.Page
}

type EnvelopeData struct {
	Envelope  *envelope.Envelope
	Tab       string
	Endpoints []endpoint.Endpoint
}

type AttachmentsData struct {
	Attachments []envelope.Attachment
	Page        paginate.Page
}

type EndpointsData struct {
	Endpoints []endpoint.Endpoint
}
