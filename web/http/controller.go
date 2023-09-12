package http

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/web"
	c "github.com/ItsNotGoodName/smtpbridge/web/components"
	"github.com/ItsNotGoodName/smtpbridge/web/meta"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
)

type head struct {
	tags []byte
}

func newHead() head {
	return head{
		tags: []byte(strings.Join(web.HeadTags, "\n")),
	}
}

func (h head) Render(ctx context.Context, w io.Writer) error {
	w.Write(h.tags)
	return nil
}

type Controller struct {
	app  core.App
	head head
	meta meta.Meta
}

func NewController(app core.App, timeHourFormat string) Controller {
	return Controller{
		app:  app,
		head: newHead(),
		meta: meta.Meta{
			TimeHourFormat: timeHourFormat,
		},
	}
}

func (ct Controller) Meta(r *http.Request) meta.Meta {
	ct.meta.Anonymous = ct.app.AuthHTTPAnonymous()
	ct.meta.Route = routes.Route(r.RequestURI)
	return ct.meta
}

func (ct Controller) Page(w http.ResponseWriter, r *http.Request, body templ.Component) {
	csrfToken := csrf.Token(r)
	c.Base(ct.head, body, csrfToken).Render(r.Context(), w)
}

func (ct Controller) Component(w http.ResponseWriter, r *http.Request, body templ.Component) {
	body.Render(r.Context(), w)
}

func (ct Controller) Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	http.Error(w, err.Error(), code)
}
