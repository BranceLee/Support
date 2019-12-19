package coreservice

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

type middleware struct {
}

// Deal with http handleFunc
func (m middleware) apply(h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		fn := middlewares[i]
		h = fn(h)
	}
	return h
}

func (m *middleware) configureSentry(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTags(map[string]string{
				"hostname":       r.Host,
				"request_host":   r.URL.Hostname(),
				"request_path":   r.URL.Path,
				"remote_address": r.RemoteAddr,
			})
		})
		h(w, r)
	}
}

func (m middleware) cors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		h(w, r)
	}
}
