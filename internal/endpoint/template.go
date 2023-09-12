package endpoint

import (
	"fmt"
	"text/template"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

type CreateFuncMap struct {
	URL string
}

func NewFuncMap(r CreateFuncMap) template.FuncMap {
	return template.FuncMap{
		"PermaLink": func(model any) string {
			switch model := model.(type) {
			case models.Envelope:
				return r.URL + routes.Envelope(model.Message.ID).URLString()
			case *models.Message:
				return r.URL + routes.Envelope(model.ID).URLString()
			case *models.Attachment:
				return r.URL + routes.AttachmentFile(model.FileName()).URLString()
			case []*models.Attachment:
				if len(model) == 0 {
					return ""
				}
				return r.URL + routes.AttachmentFile(model[0].FileName()).URLString()
			}

			panic(fmt.Sprintf("unknown model: %T", model))
		},
	}
}
