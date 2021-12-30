package router

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/dto"
)

func (s *Router) GetAttachments() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.StripPrefix(s.attachmentURI, http.FileServer(http.FS(s.a.AttachmentGetFS()))).ServeHTTP(rw, r)
	}
}

//go:embed template
var templateFS embed.FS

func (s *Router) GetIndex() http.HandlerFunc {
	type Data struct {
		Messages []dto.Message
	}

	index, err := template.ParseFS(templateFS, "template/index.html")
	if err != nil {
		log.Fatal("router.Router.GetIndex", err)
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		messages, err := s.a.MessageList(&app.MessageListRequest{AttachmentPath: s.attachmentURI, Page: 0})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("router.Router.GetIndex:", len(messages), "messages")
		err = index.Execute(rw, Data{Messages: messages})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
