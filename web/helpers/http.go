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

type Meta struct {
	Anonymous bool
	CSRF      string
}

func Render(c *fiber.Ctx, name string, data fiber.Map, layouts ...string) error {
	meta := Meta{}
	if c.Locals(AnonymousContextKey) != nil {
		meta.Anonymous = true
	}
	if csrfRaw := c.Locals(CSRFContextKey); csrfRaw != "" {
		meta.CSRF = csrfRaw.(string)
	}

	data["Meta"] = meta

	return c.Render(name, data, layouts...)
}
