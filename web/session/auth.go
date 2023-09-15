package session

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/pkg/htmx"
	"github.com/ItsNotGoodName/smtpbridge/web/routes"
	"github.com/gorilla/sessions"
)

func AuthLogin(w http.ResponseWriter, r *http.Request, ss sessions.Store, id int64) error {
	session, _ := ss.New(r, "auth")

	session.Values["id"] = id

	return ss.Save(r, w, session)
}

func AuthLogout(w http.ResponseWriter, r *http.Request, ss sessions.Store) error {
	session, _ := ss.Get(r, "auth")

	delete(session.Values, "id")

	return ss.Save(r, w, session)
}

func AuthRequire(app core.App, ss sessions.Store) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !app.AuthHTTPAnonymous() {
				session, _ := ss.Get(r, "auth")
				_, ok := session.Values["id"]
				if !ok {
					if htmx.GetRequest(r) {
						htmx.SetRedirect(w, routes.Login().String())
						return
					}

					http.Redirect(w, r, routes.Login().String(), http.StatusTemporaryRedirect)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

func AuthRestrict(app core.App, ss sessions.Store) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app.AuthHTTPAnonymous() {
				if htmx.GetRequest(r) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				http.Redirect(w, r, routes.Index().String(), http.StatusTemporaryRedirect)
				return
			}

			session, _ := ss.Get(r, "auth")
			_, ok := session.Values["id"]
			if ok {
				if htmx.GetRequest(r) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				http.Redirect(w, r, routes.Index().String(), http.StatusTemporaryRedirect)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
