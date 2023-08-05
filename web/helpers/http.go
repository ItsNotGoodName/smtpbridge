package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const CSRFContextKey = "csrf_key"
const AnonymousContextKey = "anonymous_key"

func Redirect(c *fiber.Ctx, path string) error {
	c.Set("HX-Redirect", path)
	return nil
}

func Error(c *fiber.Ctx, err error, codes ...int) error {
	c.Set("HX-Redirect", "/something-went-wrong")
	log.Err(err).Msg("Request failed")
	return nil
}

func NotFound(c *fiber.Ctx) error {
	return Render(c.Status(404), "404", fiber.Map{})
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
