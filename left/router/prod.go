//go:build !dev

package router

import (
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/go-chi/chi/v5"
)

func hookRoutes(r chi.Router) {
	mountFS(r, web.FS())
}

func hookMiddleware(r chi.Router) {}
