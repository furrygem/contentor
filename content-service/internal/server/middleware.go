package server

import (
	"net/http"

	"github.com/furrygem/contentor/content-service/pkg/webutils"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			webutils.WriteHTTPCode(w, http.StatusUnauthorized)
			return
		}
		if auth == "creds" {
			next.ServeHTTP(w, r)
			return
		}
		webutils.WriteHTTPCode(w, http.StatusForbidden)
		return
	})
}
