package router

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/left"
	"github.com/go-chi/chi/v5"
)

func handleMessageGet(w left.WebRepository, a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, "uuid")

		res, err := a.MessageGet(&app.MessageGetRequest{UUID: uuid})
		if err != nil {
			status := http.StatusInternalServerError
			if err == core.ErrMessageNotFound {
				status = http.StatusNotFound
			}
			http.Error(rw, err.Error(), status)
			return
		}

		data := left.MessageData{Message: res.Message}
		render(rw, w, left.MessagePage, &data)
	}
}

func handleMessageSendGet(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		req := &app.MessageSendRequest{
			UUID: chi.URLParam(r, "uuid"),
		}

		err := a.MessageSend(req)
		if err != nil {
			status := http.StatusInternalServerError
			if err == core.ErrMessageNotFound {
				status = http.StatusNotFound
			}
			http.Error(rw, err.Error(), status)
			return
		}

		http.Redirect(rw, r, "/message/"+req.UUID, http.StatusTemporaryRedirect)
	}
}

func handleMessageDeleteGet(a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		req := &app.MessageSendRequest{
			UUID: chi.URLParam(r, "uuid"),
		}

		err := a.MessageDelete(app.MessageDeleteRequest{UUID: req.UUID})
		if err != nil {
			status := http.StatusInternalServerError
			if err == core.ErrMessageNotFound {
				status = http.StatusNotFound
			}
			http.Error(rw, err.Error(), status)
			return
		}

		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
	}
}
