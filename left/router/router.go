package router

import (
	"context"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/left/assets"
	"github.com/ItsNotGoodName/smtpbridge/left/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

type Router struct {
	addr string
	r    chi.Router
}

func New(addr string, c *controller.Controller, dataFS fs.FS) *Router {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/assets/*", handlePrefixFS("/assets/", assets.FS()))

	r.Get("/data/*", handlePrefixFS("/data/", dataFS))

	r.Get("/", c.IndexGet)

	r.Route("/envelopes/{id}", func(r chi.Router) {
		r.Get("/", mwMultiplexAction(c.EnvelopeGet, nil, c.EnvelopeDelete))
		r.Delete("/", c.EnvelopeDelete)
		r.Get("/html", c.EnvelopeHTMLGet)
		r.Post("/send", c.EnvelopeSendPost)
	})

	r.Get("/attachments", c.AttachmentList)

	r.Route("/endpoints", func(r chi.Router) {
		r.Get("/", c.EndpointList)
		r.Post("/test", c.EndpointTestPost)
	})

	return &Router{
		addr: addr,
		r:    r,
	}
}

func (r *Router) Start() {
	log.Println("router.Router.Start: HTTP server listening on", r.addr)
	if err := http.ListenAndServe(r.addr, r.r); err != nil {
		log.Fatalln("router.Router.Start:", err)
	}
}

func (r *Router) Run(ctx context.Context, doneC chan<- struct{}) {
	go r.Start()
	doneC <- struct{}{}
}
