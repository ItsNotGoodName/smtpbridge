package helpers

import (
	"github.com/gofiber/fiber/v2"
)

const CSRFContextKey = "csrf_key"
const AnonymousContextKey = "anonymous_key"

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

	return Render(c, "something-went-wrong", fiber.Map{"Error": err})
}

func NotFound(c *fiber.Ctx) error {
	return Render(c.Status(404), "404", nil)
}

type Meta struct {
	Anonymous bool
	CSRF      string
}

func Render(c *fiber.Ctx, name string, data fiber.Map, layouts ...string) error {
	data["Meta"] = Meta{
		Anonymous: c.Locals(AnonymousContextKey) != nil,
		CSRF:      c.Locals(CSRFContextKey).(string),
	}

	return c.Render(name, data, layouts...)
}
