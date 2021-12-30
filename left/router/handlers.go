package router

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

func (s *Router) GetAttachments(prefix string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.StripPrefix(prefix, http.FileServer(http.Dir(s.attDir))).ServeHTTP(rw, r)
	}
}

//go:embed template
var templateFS embed.FS

func (s *Router) GetIndex() http.HandlerFunc {
	type Data struct {
		Messages []app.Message
	}

	index, err := template.ParseFS(templateFS, "template/index.html")
	if err != nil {
		log.Fatal("router.Router.GetIndex", err)
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		msg, err := s.messageREPO.GetMessages(10, 0)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		for i := range msg {
			msg[i].Attachments, err = s.attachmentREPO.GetAttachments(&msg[i])
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		}

		data := Data{
			Messages: msg,
		}

		index.Execute(rw, data)
	}
}
