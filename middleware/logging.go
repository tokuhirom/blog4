package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// responseWriterWrapper is a wrapper to capture HTTP response details
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log in Combined Log Format
		fmt.Printf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			ww.statusCode,
			ww.size,
			r.Referer(),
			r.UserAgent(),
		)
	})
}
