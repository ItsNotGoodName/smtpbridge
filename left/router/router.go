package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/left/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r    chi.Router
	addr string
}

func New(app dto.App, rd api.Renderer, addr string) *Router {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	hookMiddleware(r)

	r.Get("/attachment/*", mwCacheControl(handleFS("/attachment/", app.AttachmentFS()), SecondsInYear))

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/version", render(rd, api.VersionGet(app)))
		r.Get("/info", render(rd, api.InfoGet(app)))
		r.Get("/events", render(rd, api.EventsGet(app)))
		r.Get("/messages", render(rd, api.MessagesGet(app)))
		r.Get("/message/{id}", render(rd, api.MessageGet(app)))
		r.Delete("/message/{id}", render(rd, api.MessageDelete(app)))
		r.Get("/message/{id}/events", render(rd, api.MessageEventsGet(app)))
	})

	hookRoutes(r)

	return &Router{
		r:    r,
		addr: addr,
	}
}

func (r *Router) Start() {
	log.Println("router.Router.Start: HTTP server listening on", r.addr)
	err := http.ListenAndServe(r.addr, r.r)
	if err != nil {
		log.Fatalln("router.Router.Start:", err)
	}
}

func (r *Router) Run(ctx context.Context, done chan<- struct{}) {
	go r.Start()
	done <- struct{}{}
}
