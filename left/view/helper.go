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
	"pageToggleAscendingLink": func(page paginate.Page) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", page.Page, !page.Ascending, page.Limit)
	},
	"pageLimitLink": func(page paginate.Page, limit int) string {
		return fmt.Sprintf("?page=%d&ascending=%v&limit=%d", ((page.Page-1)*page.Limit)/limit+1, page.Ascending, limit)
	},
	"attachmentDataLink": func(att *envelope.Attachment) string {
		return fmt.Sprintf("/data/%s", att.FileName())
	},
	"envelopeIDLink": func(id int64) string {
		return fmt.Sprintf("/envelopes/%d", id)
	},
}
