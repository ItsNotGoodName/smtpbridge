package router

import (
	"encoding/json"
	"log"
	"net/http"
)

func renderJSON(rw http.ResponseWriter, code int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		log.Println("router.renderJSON:", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func renderError(rw http.ResponseWriter, code int, err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(struct{ Error error }{err}); err != nil {
		log.Println("router.renderError:", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
