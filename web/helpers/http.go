package helpers

import (
	"github.com/gofiber/fiber/v2"
)

func IsHTMXRequest(c *fiber.Ctx) bool {
	_, isHTMXRequest := c.GetReqHeaders()["Hx-Request"]
	return isHTMXRequest
}

func Redirect(c *fiber.Ctx, path string) error {
	if IsHTMXRequest(c) {
		c.Set("HX-Location", path)
		return nil
	} else {
		return c.Redirect(path)
	}
}

func Error(c *fiber.Ctx, err error, codes ...int) error {
	if IsHTMXRequest(c) {
		c.Set("HX-Redirect", "/something-went-wrong")
	}

	return c.Render("something-went-wrong", fiber.Map{"Error": err})
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).Render("404", nil)
}
