package http

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web/controllers"
	"github.com/ItsNotGoodName/smtpbridge/web/middleware"
	"github.com/gofiber/fiber/v2"
)

func routes(http *fiber.App, app core.App, retentionPolicy models.RetentionPolicy) {
	http.Get("/", middleware.App(app, controllers.Index(retentionPolicy)))
	http.Route("/envelopes", func(http fiber.Router) {
		http.Get("/", middleware.App(app, controllers.Envelopes))
		http.Get("/new", controllers.EnvelopesNewGet)
		http.Post("/new", middleware.App(app, controllers.EnvelopesNewPost))
		http.Route("/:id", func(http fiber.Router) {
			http.Get("/", middleware.AppID(app, controllers.EnvelopesView))
			http.Delete("/", middleware.AppID(app, controllers.EnvelopesDelete))
			http.Get("/html", middleware.AppID(app, controllers.EnvelopesViewHTML))
		})
	})
	http.Route("/attachments", func(http fiber.Router) {
		http.Get("/", middleware.App(app, controllers.Attachments))
	})
	http.Get("/endpoints", middleware.App(app, controllers.Endpoints))
	http.Route("/rules", func(http fiber.Router) {
		http.Get("/", middleware.App(app, controllers.Rules))
		http.Route("/:id", func(http fiber.Router) {
			http.Post("/enable", middleware.AppID(app, controllers.RulesEnable))
		})
	})
	http.Post("/send", middleware.App(app, controllers.Send))
	http.Post("/vacuum", middleware.App(app, controllers.Vacuum))
	http.Group("/files", controllers.Files(app))
	http.Route("/p", func(http fiber.Router) {
		http.Get("/storage-table", middleware.App(app, controllers.PStorageTable))
		http.Get("/recent-envelopes-table", middleware.App(app, controllers.PRecentEnvelopesTable))
	})
}
