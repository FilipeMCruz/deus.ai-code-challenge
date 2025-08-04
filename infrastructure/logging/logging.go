// Package logging is responsible for logging all incoming requests
package logging

import (
	"log"
	"net/http"
)

// WrapLogging wrap the handler so that all requests passed are logged
func WrapLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		handler.ServeHTTP(w, r)
	})
}
