package recovery

import (
	"log"
	"net/http"
)

// WrapPanicRecovery wrap the handler so that all requests passed, in case of panic, return a 500
func WrapPanicRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
