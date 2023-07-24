package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/procs"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type LoginData struct {
	Flash string
}

func Login(c *fiber.Ctx, cc core.Context) error {
	return c.Render("login", LoginData{})
}

func LoginPost(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	// Request
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Execute
	err := procs.AuthHTTPLogin(cc, username, password)
	if err != nil {
		if helpers.IsHTMXRequest(c) {
			return c.Render("login", LoginData{Flash: err.Error()}, "form")
		}
		return c.Render("login", LoginData{Flash: err.Error()})
	}

	// Response
	sess, err := store.Get(c)
	if err != nil {
		return helpers.Error(c, err)
	}

	sess.Set("auth", true)
	if err := sess.Save(); err != nil {
		return helpers.Error(c, err)
	}

	return helpers.Redirect(c, "/")
}

func Logout(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		return helpers.Error(c, err)
	}

	sess.Delete("auth")
	if err := sess.Save(); err != nil {
		return helpers.Error(c, err)
	}

	return helpers.Redirect(c, "/login")
}

func UserRequire(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	if procs.AuthHTTPAnonymous(cc) {
		return helpers.NotFound(c)
	}

	return authRequire(c, cc, store)
}

func AuthRequire(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	if procs.AuthHTTPAnonymous(cc) {
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
		return helpers.Redirect(c, "/login")
	}

	return c.Next()
}

func AuthRestrict(c *fiber.Ctx, cc core.Context, store *session.Store) error {
	if procs.AuthHTTPAnonymous(cc) {
		return helpers.Redirect(c, "/")
	}

	sess, err := store.Get(c)
	if err != nil {
		return helpers.Error(c, err)
	}

	auth := sess.Get("auth")
	if auth != nil {
		return helpers.Redirect(c, "/")
	}

	return c.Next()
}
