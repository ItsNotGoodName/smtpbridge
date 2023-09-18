package components

import (
	"errors"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

func (p RuleFormProps) Route() routes.Route {
	if p.Create {
		return routes.RuleCreate()
	}
	return routes.Rule(p.Data.ID)
}

func (p RuleFormProps) WithError(err error) RuleFormProps {
	var fieldErr models.FieldError
	if errors.As(err, &fieldErr) {
		if fieldErr.Field == models.FieldName {
			p.Data.NameError = fieldErr.Unwrap().Error()
		} else if fieldErr.Field == models.FieldExpression {
			p.Data.ExpressionErr = fieldErr.Unwrap()
		} else {
			p.Flash = Flash(FlashTypeError, FlashMessage(err.Error()))
		}

		return p
	}

	p.Flash = Flash(FlashTypeError, FlashMessage(err.Error()))
	return p
}
