package view

import (
	"fmt"
	"html/template"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/core/envelope"
	"github.com/ItsNotGoodName/smtpbridge/core/paginate"
)

var helperMap template.FuncMap = template.FuncMap{
	"dateFormat": func(date time.Time) string {
		return date.Format("Jan _2 2006 15:04:05")
	},
	"pageQuery": func(page paginate.Page, newPage int) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", newPage, page.Ascending, page.Limit)
	},
	"pageQueryToggleAscending": func(page paginate.Page) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", page.Page, !page.Ascending, page.Limit)
	},
	"pageQueryLimit": func(page paginate.Page, limit int) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", ((page.Page-1)*page.Limit)/limit+1, page.Ascending, limit)
	},
	"tabQuery": func(tab string) string {
		return fmt.Sprintf("?tab=%s#tab", tab)
	},
	"dataLinkAttachment": func(att *envelope.Attachment) string {
		return fmt.Sprintf("/data/%s", att.FileName())
	},
	"envelopeLink": func(env *envelope.Envelope) string {
		return fmt.Sprintf("/envelopes/%d", env.Message.ID)
	},
	"envelopeLinkID": func(id int64) string {
		return fmt.Sprintf("/envelopes/%d", id)
	},
	"envelopeHTMLLink": func(env *envelope.Envelope) string {
		return fmt.Sprintf("/envelopes/%d/html", env.Message.ID)
	},
	"envelopeSendLink": func(env *envelope.Envelope) string {
		return fmt.Sprintf("/envelopes/%d/send", env.Message.ID)
	},
	"endpointTestLinkQuery": func(end *endpoint.Endpoint) string {
		return fmt.Sprintf("/endpoints/test?endpoint=%s", end.Name)
	},
	"endpointIcon": func(end *endpoint.Endpoint) string {
		if end.Type == endpoint.TypeTelegram {
			return "fab fa-telegram fa-lg"
		}
		if end.Type == endpoint.TypeConsole {
			return "fas fa-terminal"
		}
		return "fas fa-circle-question"
	},
}
