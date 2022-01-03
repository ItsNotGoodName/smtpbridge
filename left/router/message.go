package router

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/core"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/go-chi/chi/v5"
)

func handleMessageGet(t *web.Templater, a *app.App) http.HandlerFunc {
	type Data struct {
		Message app.Message
	}

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

		t.Render(web.PageMessage, rw, Data{Message: res.Message})
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
