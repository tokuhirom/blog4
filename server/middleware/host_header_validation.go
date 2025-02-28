package middleware

import "net/http"

func HostHeader(validHost string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host != validHost && r.URL.Path != "/healthz" {
				http.Error(w, "Invalid Host header", http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
