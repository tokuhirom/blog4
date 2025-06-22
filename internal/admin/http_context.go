package admin

import (
	"context"
	"net/http"
)

type contextKey string

const (
	httpRequestKey  contextKey = "httpRequest"
	httpResponseKey contextKey = "httpResponse"
	usernameKey     contextKey = "username"
)

// HTTPContextMiddleware injects HTTP request and response into context
func HTTPContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpRequestKey, r)
		ctx = context.WithValue(ctx, httpResponseKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetHTTPRequest gets the HTTP request from context
func GetHTTPRequest(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(httpRequestKey).(*http.Request)
	return req, ok
}

// GetHTTPResponse gets the HTTP response writer from context
func GetHTTPResponse(ctx context.Context) (http.ResponseWriter, bool) {
	resp, ok := ctx.Value(httpResponseKey).(http.ResponseWriter)
	return resp, ok
}

// GetUsername gets the username from context
func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}
