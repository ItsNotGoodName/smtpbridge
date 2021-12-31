package router

import (
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	r             *chi.Mux
	a             *app.App
	attachmentURI string
}

func New(app *app.App) *Router {
	s := Router{
		r:             chi.NewRouter(),
		a:             app,
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

	s.route()

	return &s
}

func (s *Router) route() {
	s.r.Get(s.attachmentURI+"*", s.handleAttachmentsGET())
	s.r.Get("/", s.handleIndexGET())
}

func (s *Router) Start(address string) {
	log.Println("router.Router.Start: HTTP server listening on", address)
	err := http.ListenAndServe(address, s.r)
	if err != nil {
		log.Fatalln("router.Router.Start:", err)
	}
}
