// Package api is responsible for defining the http endpoint handlers provided by the service.
package api

import (
	"fmt"
	"net/http"

	"deus.ai-code-challenge/domain"
)

const (
	errInvalidPageURL     = "invalid page url: "
	errMissingFieldPrefix = "missing request field: "
	errMissingParamPrefix = "missing query param: "
	errMarshallResponse   = "unable to write response"
	errUnmarshallRequest  = "unable to read request body"
)

// Handlers returns all the service registered url and handler pairs
func Handlers(repo domain.VisitRepository) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /api/v1/unique-visitors":  buildUniqueVisitorForPageHandler(repo),
		"POST /api/v1/user-navigation": buildUserNavigationHandler(repo),
	}
}

// writeError emulates what http.Error does but uses json instead of text to represent the data
// this also ensures that all error responses follow the same structure
func writeError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, error)
}
