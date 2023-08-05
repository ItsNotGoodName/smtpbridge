package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	h "github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

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
		return h.Render(c, "login", loginData(username, err.Error()), "form")
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

func loginData(username, flash string) fiber.Map {
	return fiber.Map{
		"Username": username,
		"Flash":    flash,
	}
}
