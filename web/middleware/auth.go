package middleware

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	h "github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func UserRequire(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	if procs.AuthHTTPAnonymous(cc) {
		return h.NotFound(c)
	}

	return authRequire(c, cc, store)
}

func AuthRequire(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	if procs.AuthHTTPAnonymous(cc) {
		c.Locals(h.AnonymousContextKey, true)
		return c.Next()
	}

	return authRequire(c, cc, store)
}

func authRequire(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	auth := sess.Get("auth")
	if auth == nil {
		return h.Redirect(c, "/login")
	}

	return c.Next()
}

func AuthRestrict(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	if procs.AuthHTTPAnonymous(cc) {
		return h.Redirect(c, "/")
	}

	sess, err := store.Get(c)
	if err != nil {
		return h.Error(c, err)
	}

	auth := sess.Get("auth")
	if auth != nil {
		return h.Redirect(c, "/")
	}

	return c.Next()
}
