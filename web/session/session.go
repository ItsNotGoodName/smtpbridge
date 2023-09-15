package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func ConfigureOptions(opts *sessions.Options) {
	opts.SameSite = http.SameSiteLaxMode
	opts.HttpOnly = true
}
