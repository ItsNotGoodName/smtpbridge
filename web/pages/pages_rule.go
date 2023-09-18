package pages

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
	c "github.com/ItsNotGoodName/smtpbridge/web/components"
	"github.com/ItsNotGoodName/smtpbridge/web/forms"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"github.com/samber/lo"
)

func RuleExpressionCheck(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse
		expression := r.FormValue("expression")

		// Request
		err := app.RuleExpressionCheck(ctx, expression)

		// Render
		ct.Component(w, r, c.RuleFormExpressionLabel(c.RuleFormExpressionLabelProps{
			Err: err,
		}))
	}
}

func RuleView(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Request
		ruleEndpoints, err := app.RuleEndpointsGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}
		endpoints, err := app.EndpointList(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}
		expressionErr := app.RuleExpressionCheck(ctx, ruleEndpoints.Rule.Expression)

		// Computed
		endpointsSelections := helpers.EndpointsSelections(ruleEndpoints.Endpoints, endpoints)

		// Render
		ct.Page(w, r, ruleView(ct.Meta(r), ruleViewProps{
			Rule: ruleEndpoints.Rule,
			RuleFormProps: c.RuleFormProps{
				Data: c.RuleFormData{
					ID:                  ruleEndpoints.Rule.ID,
					Internal:            ruleEndpoints.Rule.Internal,
					Name:                ruleEndpoints.Rule.Name,
					Expression:          ruleEndpoints.Rule.Expression,
					ExpressionErr:       expressionErr,
					EndpointsSelections: endpointsSelections,
					Endpoints:           endpoints,
				},
			},
		}))
	})
}

func RuleUpdate(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Parse
		var form forms.RuleUpdate
		if err := helpers.DecodeForm(w, r, &form); err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}
		var endpointIDs []int64
		for _, v := range r.Form["endpoints"] {
			endpointID, err := strconv.Atoi(v)
			if err != nil {
				ct.Error(w, r, err, http.StatusBadRequest)
				return
			}

			endpointIDs = append(endpointIDs, int64(endpointID))
		}

		// Request
		updateErr := app.RuleUpdate(ctx, models.DTORuleUpdate{
			ID:         id,
			Name:       &form.Name,
			Expression: &form.Expression,
			Endpoints:  &form.Endpoints,
		})
		ruleEndpoints, err := app.RuleEndpointsGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}
		endpoints, err := app.EndpointList(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}
		expressionError := app.RuleExpressionCheck(ctx, form.Expression)

		// Computed
		var endpointsSelections []bool
		for _, end := range endpoints {
			selection := lo.Contains(form.Endpoints, end.ID)
			endpointsSelections = append(endpointsSelections, selection)
		}

		// Render
		props := c.RuleFormProps{
			Data: c.RuleFormData{
				ID:                  ruleEndpoints.Rule.ID,
				Internal:            ruleEndpoints.Rule.Internal,
				Name:                form.Name,
				Expression:          form.Expression,
				Endpoints:           endpoints,
				EndpointsSelections: endpointsSelections,
				ExpressionErr:       expressionError,
			},
		}
		if updateErr != nil {
			props = props.WithError(updateErr)
		} else {
			props.Flash = c.Flash(c.FlashTypeSuccess, c.FlashMessage("Updated"))
		}
		ct.Component(w, r, c.RuleForm(props))
	})
}

func RuleCreateView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Request
		endpoints, err := app.EndpointList(ctx)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		// Computed
		endpointsSelections := make([]bool, len(endpoints))

		// Render
		ct.Page(w, r, ruleCreateView(ct.Meta(r), ruleCreateViewProps{
			RuleFormProps: c.RuleFormProps{
				Create: true,
				Data: c.RuleFormData{
					Endpoints:           endpoints,
					EndpointsSelections: endpointsSelections,
				},
			},
		}))
	}
}

func RuleCreate(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse
		var form forms.RuleCreate
		if err := helpers.DecodeForm(w, r, &form); err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}

		// Request
		id, createErr := app.RuleCreate(ctx, models.DTORuleCreate{
			Name:       form.Name,
			Expression: form.Expression,
			Endpoints:  form.Endpoints,
		})
		if createErr != nil {
			// Request
			endpoints, err := app.EndpointList(ctx)
			if err != nil {
				ct.Error(w, r, err, getCode(err))
				return
			}
			expressionError := app.RuleExpressionCheck(ctx, form.Expression)

			// Computed
			var endpointsSelections []bool
			for _, end := range endpoints {
				selection := lo.Contains(form.Endpoints, end.ID)
				endpointsSelections = append(endpointsSelections, selection)
			}

			// Render
			ct.Component(w, r, c.RuleForm(c.RuleFormProps{
				Create: true,
				Data: c.RuleFormData{
					Name:                form.Name,
					Expression:          form.Expression,
					ExpressionErr:       expressionError,
					Endpoints:           endpoints,
					EndpointsSelections: endpointsSelections,
				},
			}.WithError(createErr)))
			return
		}

		htmx.SetRedirect(w, routes.Rule(id).URLString())
	}
}

func RuleDelete(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Request
		err := app.RuleDelete(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func RuleToggle(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Parse
		enable := r.FormValue("enable") == "true"

		// Request
		err := app.RuleUpdate(ctx, models.DTORuleUpdate{
			ID:     id,
			Enable: &enable,
		})
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		// Render
		ct.Component(w, r, c.RuleToggleButton(c.RuleToggleButtonProps{
			Enable: enable,
			ID:     id,
		}))
	})
}
