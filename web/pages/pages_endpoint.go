package pages

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
	c "github.com/ItsNotGoodName/smtpbridge/web/components"
	"github.com/ItsNotGoodName/smtpbridge/web/forms"
	"github.com/ItsNotGoodName/smtpbridge/web/helpers"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
)

func EndpointCreateView(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Render
		ct.Page(w, r, endpointCreate(ct.Meta(r), endpointCreateProps{
			EndpointFormProps: c.EndpointFormProps{
				Create: true,
				Data: c.EndpointFormData{
					TitleTemplate: models.EndpointTitleTemplate,
					BodyTemplate:  models.EndpointBodyTemplate,
				},
			},
		}))
	}
}

func EndpointCreate(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse
		var form forms.EndpointCreate
		if err := helpers.DecodeForm(w, r, &form); err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}
		formConfig := make(map[string]string)
		for _, config := range form.Config {
			formConfig[config.Key] = config.Value
		}

		// Request
		id, createErr := app.EndpointCreate(ctx, models.DTOEndpointCreate{
			Name:              form.Name,
			Kind:              form.Kind,
			TextDisable:       form.TextDisable,
			AttachmentDisable: form.AttachmentDisable,
			TitleTemplate:     form.TitleTemplate,
			BodyTemplate:      form.BodyTemplate,
			Config:            formConfig,
		})
		if createErr != nil {
			// Computed
			var fields []c.EndpointFormConfigField
			item := helpers.EndpointSchema().Get(form.Kind)
			for _, field := range item.Fields {
				fields = append(fields, c.EndpointFormConfigField{
					Name:        field.Name,
					Description: field.Description,
					Key:         field.Key,
					Multiline:   field.Multiline,
					Value:       formConfig[field.Key],
				})
			}

			// Render
			ct.Component(w, r, c.EndpointForm(c.EndpointFormProps{
				Create: true,
				Data: c.EndpointFormData{
					Name:              form.Name,
					Kind:              form.Kind,
					TextDisable:       form.TextDisable,
					AttachmentDisable: form.AttachmentDisable,
					TitleTemplate:     form.TitleTemplate,
					BodyTemplate:      form.BodyTemplate,
					EndpointFormConfigProps: c.EndpointFormConfigProps{
						Fields: fields,
					},
				},
			}.WithError(createErr)))
			return
		}

		htmx.SetRedirect(w, routes.Endpoint(id).URLString())
	}
}

func EndpointView(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Request
		endpoint, err := app.EndpointGet(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		// Computed
		item := helpers.EndpointSchema().Get(endpoint.Kind)
		var fields []c.EndpointFormConfigField
		for _, field := range item.Fields {
			fields = append(fields, c.EndpointFormConfigField{
				Name:        field.Name,
				Description: field.Description,
				Key:         field.Key,
				Multiline:   field.Multiline,
				Value:       endpoint.Config[field.Key],
			})
		}

		ct.Page(w, r, endpointView(ct.Meta(r), endpointViewProps{
			Endpoint: endpoint,
			EndpointFormProps: c.EndpointFormProps{
				Data: c.EndpointFormData{
					ID:                endpoint.ID,
					Internal:          endpoint.Internal,
					Name:              endpoint.Name,
					AttachmentDisable: endpoint.AttachmentDisable,
					TextDisable:       endpoint.TextDisable,
					TitleTemplate:     endpoint.TitleTemplate,
					BodyTemplate:      endpoint.BodyTemplate,
					Kind:              endpoint.Kind,
					EndpointFormConfigProps: c.EndpointFormConfigProps{
						Fields: fields,
					},
				},
			},
		}))
	})
}

func EndpointUpdate(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Parse
		var form forms.EndpointUpdate
		if err := helpers.DecodeForm(w, r, &form); err != nil {
			ct.Error(w, r, err, http.StatusBadRequest)
			return
		}
		formConfig := make(map[string]string)
		for _, config := range form.Config {
			formConfig[config.Key] = config.Value
		}

		// Request
		updateErr := app.EndpointUpdate(ctx, models.DTOEndpointUpdate{
			ID:                id,
			Name:              &form.Name,
			Kind:              &form.Kind,
			TextDisable:       &form.TextDisable,
			AttachmentDisable: &form.AttachmentDisable,
			TitleTemplate:     &form.TitleTemplate,
			BodyTemplate:      &form.BodyTemplate,
			Config:            &formConfig,
		})

		// Computed
		item := helpers.EndpointSchema().Get(form.Kind)
		var fields []c.EndpointFormConfigField
		for _, field := range item.Fields {
			fields = append(fields, c.EndpointFormConfigField{
				Name:        field.Name,
				Description: field.Description,
				Key:         field.Key,
				Multiline:   field.Multiline,
				Value:       formConfig[field.Key],
			})
		}

		// Render
		props := c.EndpointFormProps{
			Data: c.EndpointFormData{
				ID:                id,
				Name:              form.Name,
				AttachmentDisable: form.AttachmentDisable,
				TextDisable:       form.TextDisable,
				TitleTemplate:     form.TitleTemplate,
				BodyTemplate:      form.BodyTemplate,
				Kind:              form.Kind,
				EndpointFormConfigProps: c.EndpointFormConfigProps{
					Fields: fields,
				},
			},
		}
		if updateErr != nil {
			props = props.WithError(updateErr)
		} else {
			props.Flash = c.Flash(c.FlashTypeSuccess, c.FlashMessage("Updated."))
		}
		ct.Component(w, r, c.EndpointForm(props))
	})
}

func EndpointDelete(ct Controller, app core.App) http.HandlerFunc {
	return withID(ct, func(w http.ResponseWriter, r *http.Request, id int64) {
		ctx := r.Context()

		// Request
		err := app.EndpointDelete(ctx, id)
		if err != nil {
			ct.Error(w, r, err, getCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func EndpointFormConfigComponent(ct Controller, app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse
		kind := r.URL.Query().Get("Kind")
		item := helpers.EndpointSchema().Get(kind)

		// Computed
		var fields []c.EndpointFormConfigField
		for _, field := range item.Fields {
			fields = append(fields, c.EndpointFormConfigField{
				Name:        field.Name,
				Description: field.Description,
				Key:         field.Key,
				Multiline:   field.Multiline,
			})
		}

		ct.Page(w, r, c.EndpointFormConfig(c.EndpointFormConfigProps{
			Fields: fields,
		}))
	}
}
