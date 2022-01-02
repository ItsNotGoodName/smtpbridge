package router

import (
	_ "embed"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/dto"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/ItsNotGoodName/smtpbridge/pkg/paginate"
)

func (s *Router) handleAttachmentsGET() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.StripPrefix(s.attachmentURI, http.FileServer(http.FS(s.a.AttachmentGetFS()))).ServeHTTP(rw, r)
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
