package components

import (
	"errors"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

func (p EndpointFormProps) Route() routes.Route {
	if p.Create {
		return routes.EndpointCreate()
	}
	return routes.Endpoint(p.Data.ID)
}

func (p EndpointFormProps) WithError(err error) EndpointFormProps {
	var fieldErr models.FieldError
	if errors.As(err, &fieldErr) {
		if fieldErr.Field == models.FieldName {
			p.Data.NameError = fieldErr.Unwrap().Error()
		} else if fieldErr.Field == models.FieldTitleTemplate {
			p.Data.TitleTemplateError = fieldErr.Unwrap().Error()
		} else if fieldErr.Field == models.FieldBodyTemplate {
			p.Data.BodyTemplateError = fieldErr.Unwrap().Error()
		} else if fieldErr.Field == models.FieldKind {
			p.Data.KindError = fieldErr.Unwrap().Error()
		} else if fieldErr.Field == models.FieldConfig {
			p.Data.EndpointFormConfigProps.Error = fieldErr.Unwrap().Error()
		} else {
			p.Flash = Flash(FlashTypeError, FlashMessage(err.Error()))
		}

		return p
	}

	p.Flash = Flash(FlashTypeError, FlashMessage(err.Error()))
	return p
}
