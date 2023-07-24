package helpers

import (
	"encoding/json"
	"html/template"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/build"
	"github.com/ItsNotGoodName/smtpbridge/web"
	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
)

var Map template.FuncMap = template.FuncMap{
	"build": func() build.Build {
		return build.Current
	},
	"headTags": func() template.HTML {
		return template.HTML(web.HeadTags)
	},
	"timeFormat": func(date time.Time) string {
		return date.Local().Format("Jan _2 2006 15:04:05")
	},
	"timeHumanize": func(date time.Time) string {
		return humanize.Time(date)
	},
	"bytesHumanize": func(bytes int64) string {
		return humanize.Bytes(uint64(bytes))
	},
	"json": func(data any) string {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}

		return string(jsonData)
	},
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"query": func(queries map[string]string, vals ...any) template.URL {
		return template.URL(Query(queries, vals...))
	},
}

func IsHTMXRequest(c *fiber.Ctx) bool {
	_, isHTMXRequest := c.GetReqHeaders()["Hx-Request"]
	return isHTMXRequest
}
