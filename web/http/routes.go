package http

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/web/controllers"
	"github.com/ItsNotGoodName/smtpbridge/web/inject"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func route(app core.App, store *session.Store, http fiber.Router) {
	userRequire := inject.AppStore(app, store, controllers.UserRequire)
	authRequire := inject.AppStore(app, store, controllers.AuthRequire)
	authSkip := inject.AppStore(app, store, controllers.AuthRestrict)

	http.Get("/login", authSkip, inject.App(app, controllers.Login))
	http.Post("/login", authSkip, inject.AppStore(app, store, controllers.LoginPost))
	http.Post("/logout", userRequire, inject.AppStore(app, store, controllers.Logout))

	http.Get("/", authRequire, inject.App(app, controllers.Index))
	http.Route("/index", func(http fiber.Router) {
		http.Get("/storage-table", authRequire, inject.App(app, controllers.IndexStorageTable))
		http.Get("/recent-envelopes-table", authRequire, inject.App(app, controllers.IndexRecentEnvelopesTable))
	})

	http.Route("/envelopes", func(http fiber.Router) {
		http.Get("/", authRequire, inject.App(app, controllers.Envelopes))
		http.Delete("/", authRequire, inject.App(app, controllers.EnvelopesDelete))
		http.Get("/new", authRequire, controllers.EnvelopeNew)
		http.Post("/new", authRequire, inject.App(app, controllers.EnvelopeNewPost))
		http.Route("/:id", func(http fiber.Router) {
			http.Get("/", authRequire, inject.AppID(app, controllers.Envelope))
			http.Delete("/", authRequire, inject.AppID(app, controllers.EnvelopeDelete))
			http.Get("/html", authRequire, inject.AppID(app, controllers.EnvelopeHTML))
		})
	})

	http.Route("/attachments", func(http fiber.Router) {
		http.Get("/", authRequire, inject.App(app, controllers.Attachments))
	})

	http.Route("/endpoints", func(http fiber.Router) {
		http.Get("/", authRequire, inject.App(app, controllers.Endpoints))
		http.Route("/:id", func(http fiber.Router) {
			http.Post("/test", authRequire, inject.AppID(app, controllers.EndpointTest))
		})
	})

	http.Route("/rules", func(http fiber.Router) {
		http.Get("/", authRequire, inject.App(app, controllers.Rules))
		http.Route("/:id", func(http fiber.Router) {
			http.Post("/enable", authRequire, inject.AppID(app, controllers.RuleEnable))
		})
	})

	http.Post("/send", authRequire, inject.App(app, controllers.Send))
	http.Post("/vacuum", authRequire, inject.App(app, controllers.Vacuum))
	http.Post("/trim", authRequire, inject.App(app, controllers.Trim))
	http.Group("/files", authRequire, controllers.Files(app))

	http.Get("/something-went-wrong", controllers.SomethingWentWrong)
}
