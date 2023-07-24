package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	h "github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func loginData(flash string) fiber.Map {
	return fiber.Map{
		"Flash": flash,
	}
}

func Login(c *fiber.Ctx, cc core.Context) error {
	return h.Render(c, "login", fiber.Map{})
}

func LoginPost(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	// Request
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Execute
	err := procs.AuthHTTPLogin(cc, username, password)
	if err != nil {
		if h.IsHTMXRequest(c) {
			return h.Render(c, "login", loginData(err.Error()), "form")
		}
		return h.Render(c, "login", loginData(err.Error()))
	}

	// Response
	sess, err := store.Get(c)
	if err != nil {
		return h.Error(c, err)
	}

	sess.Set("auth", true)
	if err := sess.Save(); err != nil {
		return h.Error(c, err)
	}

	return h.Redirect(c, "/")
}

func Logout(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		return h.Error(c, err)
	}

	sess.Delete("auth")
	if err := sess.Save(); err != nil {
		return h.Error(c, err)
	}

	return h.Redirect(c, "/login")
}

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
