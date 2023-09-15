package http

import (
	"compress/flate"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/pkg/chiext"
	"github.com/ItsNotGoodName/smtpbridge/web"
	"github.com/ItsNotGoodName/smtpbridge/web/pages"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"github.com/ItsNotGoodName/smtpbridge/web/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

const paramID = "{id}"

func NewRouter(ct pages.Controller, app core.App, fileFS fs.FS, csrfSecret []byte, ss sessions.Store) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(csrf.Protect(csrfSecret, csrf.Secure(false), csrf.Path("/")))

	r.Group(func(r chi.Router) {
		r.Use(web.CacheControl)
		r.Use(middleware.NewCompressor(flate.DefaultCompression, "application/javascript", "text/css").Handler)

		chiext.MountFS(r, web.FS)
	})

	// Unauthorized
	r.Group(func(r chi.Router) {
		r.Use(session.AuthRestrict(app, ss))

		// Login
		r.Get(routes.Login().String(),
			pages.LoginView(ct, app))
		r.Post(routes.Login().String(),
			pages.Login(ct, app, ss))
	})

	// Authorized
	r.Group(func(r chi.Router) {
		r.Use(session.AuthRequire(app, ss))

		// Logout
		r.Delete(routes.Logout().String(),
			pages.Logout(ct, app, ss))

		// Index
		r.Get(routes.Index().String(),
			pages.IndexView(ct, app))

		// Envelope
		r.Get(routes.EnvelopeCreate().String(),
			pages.EnvelopeCreateView(ct, app))
		r.Post(routes.EnvelopeCreate().String(),
			pages.EnvelopeCreate(ct, app))
		r.Get(routes.Envelope(paramID).String(),
			pages.EnvelopeView(ct, app))
		r.Delete(routes.Envelope(paramID).String(),
			pages.EnvelopeDelete(ct, app))
		r.Get(routes.EnvelopeHTML(paramID).String(),
			pages.EnvelopeHTMLView(ct, app))
		r.Post(routes.EnvelopeEndpointSend(paramID).String(),
			pages.EnvelopeEndpointSend(ct, app))
		{
			envelopeListView := pages.EnvelopeListView(ct, app)
			r.Get(routes.EnvelopeList().String(),
				envelopeListView)
			r.Delete(routes.EnvelopeList().String(),
				pages.EnvelopeListDrop(ct, app, envelopeListView))
		}

		// Attachment
		r.Get(routes.AttachmentFile("*").String(),
			pages.Files(ct, app, fileFS))
		{
			attachmentListView := pages.AttachmentListView(ct, app)
			r.Get(routes.AttachmentList().String(),
				attachmentListView)
			r.Post(routes.AttachmentTrim().String(),
				pages.AttachmentTrim(ct, app, attachmentListView))
		}

		// Endpoint
		r.Get(routes.EndpointList().String(),
			pages.EndpointListView(ct, app))
		r.Post(routes.EndpointTest(paramID).String(),
			pages.EndpointTest(ct, app))

		// Traces
		{
			traceListView := pages.TraceListView(ct, app)
			r.Get(routes.TraceList().String(),
				traceListView)
			r.Delete(routes.TraceList().String(),
				pages.TraceListDrop(ct, app, traceListView))
		}

		// Rules
		r.Get(routes.RuleList().String(),
			pages.RuleListView(ct, app))
		r.Get(routes.Rule(paramID).String(),
			pages.RuleView(ct, app))
		r.Delete(routes.Rule(paramID).String(),
			pages.RuleDelete(ct, app))
		r.Post(routes.Rule(paramID).String(),
			pages.RuleUpdate(ct, app))
		r.Post(routes.RuleExpressionCheck().String(),
			pages.RuleExpressionCheck(ct, app))
		r.Get(routes.RuleCreate().String(),
			pages.RuleCreateView(ct, app))
		r.Post(routes.RuleCreate().String(),
			pages.RuleCreate(ct, app))
		r.Post(routes.RuleToggle(paramID).String(),
			pages.RuleToggle(ct, app))

		// Retention Policy
		r.Post(routes.RetentionPolicyRun().String(),
			pages.RetentionPolicyRun(ct, app))

		// Components
		r.Get(routes.StorageStatsComponent().String(),
			pages.StorageStatsComponent(ct, app))
		r.Get(routes.EnvelopeTabComponent(paramID, routes.EnvelopeTabText).String(),
			pages.EnvelopeTabComponent(ct, app, routes.EnvelopeTabText))
		r.Get(routes.EnvelopeTabComponent(paramID, routes.EnvelopeTabHTML).String(),
			pages.EnvelopeTabComponent(ct, app, routes.EnvelopeTabHTML))
		r.Get(routes.EnvelopeTabComponent(paramID, routes.EnvelopeTabAttachments).String(),
			pages.EnvelopeTabComponent(ct, app, routes.EnvelopeTabAttachments))
		r.Get(routes.RecentEnvelopeListComponent().String(),
			pages.RecentEnvelopeListComponent(ct, app))
		r.Get(routes.NullComponent().String(),
			pages.NullComponent())
	})

	return r
}
