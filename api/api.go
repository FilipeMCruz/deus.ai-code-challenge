// Package api is responsible for defining the http endpoint handlers provided by the service.
package api

import (
	"deus.ai-code-challenge/domain"
	"deus.ai-code-challenge/service"
	"errors"
	"fmt"
	"net/http"
)

const (
	errMissingFieldPrefix = "missing request field: "
	errMissingParamPrefix = "missing query param: "
	errMarshallResponse   = "unable to write response"
	errUnmarshallRequest  = "unable to read request body"
)

// Handlers returns all the service registered url and handler pairs
func Handlers(visits domain.VisitsRepository, pages domain.PageRepository) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /api/v1/unique-visitors":  buildUniqueVisitorForPageHandler(visits, pages),
		"POST /api/v1/user-navigation": buildUserNavigationHandler(visits, pages),
	}
}

// writeError emulates what http.Error does but uses json instead of text to represent the data
// this also ensures that all error responses follow the same structure
func writeError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, `{"error":"%s"}`, error)
}

// getStatusCode maps errors to status codes
func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, service.ErrPageNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
