package json

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/left/api"
)

type JSON struct{}

type Response struct {
	OK    bool        `json:"ok"`
	Code  int         `json:"code"`
	Error *Error      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}

func New() *JSON {
	return &JSON{}
}

func (JSON) Render(w http.ResponseWriter, r api.Response) {
	if r.Error != nil {
		renderError(w, r.Code, r.Error)
		return
	}
	renderJSON(w, r.Code, r.Data)
}

func renderError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(Response{
		OK:   false,
		Code: code,
		Error: &Error{
			Message: err.Error(),
		},
	}); err != nil {
		log.Println("json.renderError:", err)
	}
}

func renderJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(Response{
		OK:   true,
		Code: code,
		Data: data,
	}); err != nil {
		log.Println("json.renderJSON:", err)
	}
}
