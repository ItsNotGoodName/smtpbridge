package view

import (
	"fmt"
	"html/template"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

var helperMap template.FuncMap = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("Jan _2 2006 15:04:05")
	},
	"pageLink": func(page paginate.Page, newPage int) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", newPage, page.Ascending, page.Limit)
	},
	"toggleAscendingLink": func(page paginate.Page) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", page.Page, !page.Ascending, page.Limit)
	},
	"dataLink": func(att *envelope.Attachment) string {
		return fmt.Sprintf("/data/%s", att.FileName())
	},
	"envelopeLink": func(env *envelope.Envelope) string {
		return fmt.Sprintf("/envelope/%d", env.Message.ID)
	},
}
