package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/furrygem/contentor/content-service/pkg/webutils"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})

}

type contextKey int

const (
	keyTokenSubject     contextKey = iota
	keyTokenAllowedKeys contextKey = iota
	keyTokenOwner       contextKey = iota
	keyTokenGuest       contextKey = iota
	keyTokenShared      contextKey = iota
)

func (s *Server) authenticationMiddleware(next http.Handler) http.Handler {
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
		guest := claims["guest"]
		if guest == true {
			owner := claims["owner"]
			shared := claims["shared"]
			allowedKeys := claims["allowedKeys"]
			if owner == nil || shared == nil || allowedKeys == nil {
				webutils.WriteHTTPCodeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad token"})
				return
			}
			ctx := context.WithValue(r.Context(), keyTokenAllowedKeys, allowedKeys)
			ctx = context.WithValue(ctx, keyTokenOwner, owner)
			ctx = context.WithValue(ctx, keyTokenGuest, guest)
			ctx = context.WithValue(ctx, keyTokenShared, shared)
			r = r.Clone(ctx)
			next.ServeHTTP(w, r)
			return
		}
		user := claims["sub"].(string)
		if claims.Valid() == nil {
			// r.Header.Set("X-Auth-Token-Subject", user)
			ctx := context.WithValue(r.Context(), keyTokenSubject, user)
			r = r.Clone(ctx)
			next.ServeHTTP(w, r)
			return
		}
		webutils.WriteHTTPCode(w, http.StatusForbidden)
		return
	})
}
