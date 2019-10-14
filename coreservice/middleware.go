package coreservice

import (
	"net/http"
	"os"
)

type middleware struct {
}

// Deal with http handleFunc
func (m middleware) apply(h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := 0; i < len(middlewares); i++ {
		fn := middlewares[i]
		h = fn(h)
	}
	return h
}

func (m middleware) cors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env, _ := os.LookupEnv("ENV")
		if env == "DEV" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		h(w, r)
	}
}