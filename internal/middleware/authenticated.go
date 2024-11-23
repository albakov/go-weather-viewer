package middleware

import (
	"context"
	"github.com/albakov/go-weather-viewer/internal/controller"
	"github.com/albakov/go-weather-viewer/internal/services/session"
	"net/http"
	"strings"
)

func SetAuthenticatedUser(next http.Handler, sessionService *session.Service, excludedPaths []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range excludedPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)

				return
			}
		}

		userItem := sessionService.AuthenticatedUser(r)

		if userItem.Id != 0 {
			r = r.WithContext(context.WithValue(r.Context(), "user", userItem))
		}

		next.ServeHTTP(w, r)
	})
}

func Authenticated(next http.Handler, sessionService *session.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sessionService.IsAuthenticated(r) {
			controller.New().RedirectTo(w, "/")

			return
		}

		next.ServeHTTP(w, r)
	})
}

func NotAuthenticated(next http.Handler, sessionService *session.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionService.IsAuthenticated(r) {
			controller.New().RedirectTo(w, "/")

			return
		}

		next.ServeHTTP(w, r)
	})
}
