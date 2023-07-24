package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Flash struct {
	Flash string
}

func Login(c *fiber.Ctx, cc core.Context) error {
	return c.Render("login", Flash{})
}

func AuthLogin(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	// Request
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Execute
	err := procs.HTTPLogin(cc, username, password)
	if err != nil {
		if helpers.IsHTMXRequest(c) {
			return c.Render("login", Flash{Flash: err.Error()}, "form")
		}
		return c.Render("login", Flash{Flash: err.Error()})
	}

	// Response
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Set("auth", true)
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return redirect(c, "/")
}

func AuthLogout(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Delete("auth")
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return redirect(c, "/login")
}

func AuthRequire(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	auth := sess.Get("auth")
	if auth == nil {
		return redirect(c, "/login")
	}

	return c.Next()
}

func AuthSkip(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	auth := sess.Get("auth")
	if auth != nil {
		return redirect(c, "/")
	}

	return c.Next()
}
