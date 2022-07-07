package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/furrygem/contentor/content-service/pkg/webutils"
)

func (s *Server) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			webutils.WriteHTTPCode(w, http.StatusUnauthorized)
			return
		}
		authSplit := strings.Split(auth, " ")
		if len(authSplit) != 2 || strings.ToLower(authSplit[0]) != "bearer" {
			webutils.WriteHTTPCode(w, http.StatusUnauthorized)
			return
		}
		rawToken := authSplit[1]
		_, claims, err := s.authHandler.ParseJWT(rawToken)
		if err != nil {
			webutils.WriteHTTPCodeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		user := claims["sub"].(string)
		if claims.Valid() == nil {
			// r.Header.Set("X-Auth-Token-Subject", user)
			ctx := context.WithValue(r.Context(), "token-subject", user)
			r = r.Clone(ctx)
			next.ServeHTTP(w, r)
			return
		}
		webutils.WriteHTTPCode(w, http.StatusForbidden)
		return
	})
}
