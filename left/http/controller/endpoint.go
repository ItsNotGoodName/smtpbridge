package controller

import (
	"errors"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/core/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/left/http/view"
)

func EndpointList(endpointService endpoint.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.Render(w, http.StatusOK, view.EndpointsPage, view.EndpointsData{Endpoints: endpointService.ListEndpoint()})
	}
}

func EndpointTestPost(endpointService endpoint.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointName := r.URL.Query().Get("endpoint")
		end, err := endpointService.GetEndpoint(endpointName)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, core.ErrEndpointNotFound) {
				code = http.StatusNotFound
			}
			http.Error(w, err.Error(), code)
			return
		}

		if err := end.SendRaw(r.Context(), end.Name+" test", []endpoint.Attachment{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
