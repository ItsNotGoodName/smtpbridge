package left

import (
	"html/template"
	"io/fs"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/pkg/paginate"
)

const (
	IndexPage   Page = "index.html"
	MessagePage Page = "message.html"
	InfoPage    Page = "info.html"
)

type (
	IndexData struct {
		Messages []app.Message
		Paginate paginate.Paginate
	}

	MessageData struct {
		Message app.Message
	}

	InfoData struct {
		Info app.Info
	}

	WebRepository interface {
		GetTemplate(page Page) *template.Template
		GetAssetFS() fs.FS
	}

	Page string
)
