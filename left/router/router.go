package router

import (
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(a *app.App, t *web.Templater) http.Handler {
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

	r.Get("/attachments/*", mwCacheControl(handleFS("/attachments/", a.AttachmentGetFS())))
	r.Get("/message/{uuid}", handleMessageGet(t, a))
	r.Get("/message/{uuid}/send", handleMessageSendGet(a))
	r.Get("/assets/*", handleFS("/assets/", web.GetAssetFS()))
	r.Get("/", handleIndexGet(t, a))

	return r
}

func Start(cfg *config.Config, r http.Handler) {
	log.Println("router.Start: HTTP server listening on", cfg.HTTP.Addr)
	err := http.ListenAndServe(cfg.HTTP.Addr, r)
	if err != nil {
		log.Fatalln("router.Start:", err)
	}
}
