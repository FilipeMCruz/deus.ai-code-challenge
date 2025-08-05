// Package content is responsible for setting the content type header to application/json
package content

import (
	"net/http"
)

// WrapJsonContentType wraps the handler so that all requests reply with a response that contains the header
// Content-Type: application/json
func WrapJsonContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		handler.ServeHTTP(w, r)
	})
}
