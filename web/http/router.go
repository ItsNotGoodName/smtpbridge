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
		r.Get(routes.Envelope(pages.ParamID).String(),
			pages.EnvelopeView(ct, app))
		r.Delete(routes.Envelope(pages.ParamID).String(),
			pages.EnvelopeDelete(ct, app))
		r.Get(routes.EnvelopeHTML(pages.ParamID).String(),
			pages.EnvelopeHTMLView(ct, app))
		r.Post(routes.EnvelopeEndpointSend(pages.ParamID).String(),
			pages.EnvelopeEndpointSend(ct, app))
		{
			view := pages.EnvelopeListView(ct, app)
			r.Get(routes.EnvelopeList().String(),
				view)
			r.Delete(routes.EnvelopeList().String(),
				pages.EnvelopeListDrop(ct, app, view))
		}

		// Attachment
		r.Get(routes.AttachmentFile("*").String(),
			pages.Files(ct, app, fileFS))
		{
			view := pages.AttachmentListView(ct, app)
			r.Get(routes.AttachmentList().String(),
				view)
			r.Post(routes.AttachmentTrim().String(),
				pages.AttachmentTrim(ct, app, view))
		}

		// Endpoint
		r.Get(routes.Endpoint(pages.ParamID).String(),
			pages.EndpointView(ct, app))
		r.Post(routes.Endpoint(pages.ParamID).String(),
			pages.EndpointUpdate(ct, app))
		r.Get(routes.EndpointList().String(),
			pages.EndpointListView(ct, app))
		r.Post(routes.EndpointTest(pages.ParamID).String(),
			pages.EndpointTest(ct, app))
		r.Get(routes.EndpointCreate().String(),
			pages.EndpointCreateView(ct, app))
		r.Post(routes.EndpointCreate().String(),
			pages.EndpointCreate(ct, app))
		r.Delete(routes.Endpoint(pages.ParamID).String(),
			pages.EndpointDelete(ct, app))

		// Traces
		{
			view := pages.TraceListView(ct, app)
			r.Get(routes.TraceList().String(),
				view)
			r.Delete(routes.TraceList().String(),
				pages.TraceListDrop(ct, app, view))
		}

		// Rules
		r.Get(routes.RuleList().String(),
			pages.RuleListView(ct, app))
		r.Get(routes.Rule(pages.ParamID).String(),
			pages.RuleView(ct, app))
		r.Delete(routes.Rule(pages.ParamID).String(),
			pages.RuleDelete(ct, app))
		r.Post(routes.Rule(pages.ParamID).String(),
			pages.RuleUpdate(ct, app))
		r.Post(routes.RuleExpressionCheck().String(),
			pages.RuleExpressionCheck(ct, app))
		r.Get(routes.RuleCreate().String(),
			pages.RuleCreateView(ct, app))
		r.Post(routes.RuleCreate().String(),
			pages.RuleCreate(ct, app))
		r.Post(routes.RuleToggle(pages.ParamID).String(),
			pages.RuleToggle(ct, app))

		// Retention Policy
		r.Post(routes.RetentionPolicyRun().String(),
			pages.RetentionPolicyRun(ct, app))

		// Components
		r.Get(routes.EndpointFormConfigComponent().String(),
			pages.EndpointFormConfigComponent(ct, app))
		r.Get(routes.StorageStatsComponent().String(),
			pages.StorageStatsComponent(ct, app))
		r.Get(routes.EnvelopeTabComponent(pages.ParamID, routes.EnvelopeTabText).String(),
			pages.EnvelopeTabComponent(ct, app, routes.EnvelopeTabText))
		r.Get(routes.EnvelopeTabComponent(pages.ParamID, routes.EnvelopeTabHTML).String(),
			pages.EnvelopeTabComponent(ct, app, routes.EnvelopeTabHTML))
		r.Get(routes.EnvelopeTabComponent(pages.ParamID, routes.EnvelopeTabAttachments).String(),
			pages.EnvelopeTabComponent(ct, app, routes.EnvelopeTabAttachments))
		r.Get(routes.RecentEnvelopeListComponent().String(),
			pages.RecentEnvelopeListComponent(ct, app))
		r.Get(routes.NullComponent().String(),
			pages.NullComponent())
	})

	return r
}
