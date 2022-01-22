package controllers

import (
	"net/http"

	"github.com/FranciscoMendes10866/queues/helpers"
)

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-dango-manga-key") != helpers.AllowedHeaderValue {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
