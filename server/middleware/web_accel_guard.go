package middleware

import (
	"log"
	"net/http"
)

func CheckWebAccelHeader(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotToken := r.Header.Get("X-WebAccel-Guard")
			if gotToken != token && r.URL.Path != "/healthz" {
				log.Printf("invalid X-WebAccel-Guard header: '%s'", gotToken)
				http.Error(w, "Invalid X-WebAccel-Guard header", http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
