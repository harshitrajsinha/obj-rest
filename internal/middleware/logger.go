// Package middleware defines different middlewares around request-response cycle
package middleware

import (
	"log"
	"net/http"
	"time"
)

// CustomResponseWriter implements ResponseWriter interface to override WriteHeader()
type CustomResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

// WriteHeader overrides built-in method
func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware defines the logging middleware for request and response cycle
func LoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		crw := &CustomResponseWriter{
			statusCode:     200,
			ResponseWriter: w,
		}

		next.ServeHTTP(crw, r)

		elapsedTime := time.Since(startTime).Round(time.Millisecond)
		statuscode := crw.statusCode
		level := "[INFO]"
		if statuscode >= 400 {
			level = "[ERROR]"
		}
		log.Printf("%s %s %s %d %v", level, r.Method, r.URL.Path, statuscode, elapsedTime)
	})
}
