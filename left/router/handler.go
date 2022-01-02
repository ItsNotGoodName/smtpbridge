package router

import (
	_ "embed"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/dto"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/ItsNotGoodName/smtpbridge/pkg/paginate"
)

func handleImage(prefix string, dirFS fs.FS) http.HandlerFunc {
	h := http.StripPrefix(prefix, http.FileServer(http.FS(dirFS)))

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "max-age=31536000")
		h.ServeHTTP(rw, r)
	}
}

func (s *Router) handleIndexGET() http.HandlerFunc {
	type Data struct {
		Messages []dto.Message
		Paginate paginate.Paginate
	}

	param := "page"

	return func(rw http.ResponseWriter, r *http.Request) {
		var page int
		if p := r.URL.Query().Get(param); p != "" {
			if i, err := strconv.Atoi(p); err == nil {
				page = i
			}
		}

		res, err := s.a.MessageList(&app.MessageListRequest{Page: page})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		pag := paginate.New(*r.URL, param, res.PageMin, res.Page, res.PageMax)
		data := Data{Messages: res.Messages, Paginate: pag}

		s.t.Render(web.PageIndex, rw, data)
	}
}
