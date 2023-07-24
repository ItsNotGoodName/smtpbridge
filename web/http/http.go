package http

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/web"
	h "github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog/log"
)

type HTTP struct {
	http     *fiber.App
	shutdown context.CancelFunc
	address  string
}

func New(app core.App, shutdown context.CancelFunc, address string, bodyLimit int) HTTP {
	store := session.New()

	// Fiber
	views := web.Engine()
	views.AddFuncMap(h.Map)
	http := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 views,
		ViewsLayout:           "layouts/index",
		BodyLimit:             bodyLimit,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if h.IsHTMXRequest(c) {
				c.Set("HX-Redirect", "/something-went-wrong")
			}

			return fiber.DefaultErrorHandler(c, err)
		},
	})

	// Middleware
	http.Use(recover.New())
	http.Use(logger.New())
	http.Use(csrf.New(csrf.Config{
		ContextKey: h.CSRFContextKey,
		Extractor:  csrfExtractor(),
	}))
	web.UseAssets(http)

	route(app,
		store,
		http,
	)

	http.Use(h.NotFound)

	return HTTP{
		http:     http,
		address:  address,
		shutdown: shutdown,
	}
}

func (h HTTP) Background(ctx context.Context, doneC chan<- struct{}) {
	log.Info().Msg("Starting HTTP server on " + h.address)

	go func() {
		<-ctx.Done()
		log.Info().Msg("Gracefully shutting down HTTP server...")
		if err := h.http.ShutdownWithTimeout(5 * time.Second); err != nil {
			log.Err(err).Msg("Failed to shutdown HTTP server")
		}
		doneC <- struct{}{}
	}()

	if err := h.http.Listen(h.address); err != nil {
		log.Err(err).Msg("HTTP server failed to listen")
		h.shutdown()
	}
}
