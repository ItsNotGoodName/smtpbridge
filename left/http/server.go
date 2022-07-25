package http

import (
	"context"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/left/http/asset"
	c "github.com/ItsNotGoodName/smtpbridge/left/http/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

type Server struct {
	addr string
	r    chi.Router
}

func New(addr string, dataFS fs.FS, envelopeService envelope.Service, endpointService endpoint.Service) *Server {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/assets/*", handlePrefixFS("/assets/", asset.FS()))

	r.Get("/data/*", handlePrefixFS("/data/", dataFS))

	r.Get("/", c.IndexGet(envelopeService))

	r.Route("/envelopes/{id}", func(r chi.Router) {
		r.Get("/", mwMultiplexAction(c.EnvelopeGet(envelopeService, endpointService), nil, c.EnvelopeDelete(envelopeService)))
		r.Delete("/", c.EnvelopeDelete(envelopeService))
		r.Get("/html", c.EnvelopeHTMLGet(envelopeService))
		r.Post("/send", c.EnvelopeSendPost(envelopeService, endpointService))
	})

	r.Get("/attachments", c.AttachmentList(envelopeService))

	r.Route("/endpoints", func(r chi.Router) {
		r.Get("/", c.EndpointList(endpointService))
		r.Post("/test", c.EndpointTestPost(endpointService))
	})

	return &Server{
		addr: addr,
		r:    r,
	}
}

func (s *Server) Start() (*http.Server, <-chan struct{}) {
	ch := make(chan struct{})
	server := &http.Server{Addr: s.addr, Handler: s.r}

	go func() {
		defer close(ch)
		log.Println("http.Server.Start: HTTP server listening on", s.addr)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("http.Server.Start:", err)
		}
	}()

	return server, ch
}

func (s *Server) Run(ctx context.Context, doneC chan<- struct{}) {
	srv, ch := s.Start()               // Start HTTP server
	<-ctx.Done()                       // Wait for context
	srv.Shutdown(context.Background()) // Shutdown HTTP server
	<-ch                               // Wait for HTTP server shutodwn
	doneC <- struct{}{}                // Done
}
