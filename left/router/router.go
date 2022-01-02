package router

import (
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	addr          string
	r             *chi.Mux
	a             *app.App
	t             *web.Templater
	attachmentURI string
}

func New(cfg *config.Config, app *app.App, templater *web.Templater) *Router {
	s := Router{
		addr:          cfg.HTTP.Addr,
		r:             chi.NewRouter(),
		a:             app,
		t:             templater,
		attachmentURI: "/attachments/",
	}

	// A good base middleware stack
	s.r.Use(middleware.RequestID)
	s.r.Use(middleware.RealIP)
	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	s.r.Use(middleware.Timeout(60 * time.Second))

	s.r.Get(s.attachmentURI+"*", s.handleAttachmentsGET())
	s.r.Get("/", s.handleIndexGET())

	return &s
}

func (s *Router) Start() {
	log.Println("router.Router.Start: HTTP server listening on", s.addr)
	err := http.ListenAndServe(s.addr, s.r)
	if err != nil {
		log.Fatalln("router.Router.Start:", err)
	}
}
