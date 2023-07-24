package controllers

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/gofiber/fiber/v2"
)

func AuthLogin(c *fiber.Ctx, cc core.Context) error {
	// Request

	// Execute

	// Response
	return c.Render("login", fiber.Map{})
}
