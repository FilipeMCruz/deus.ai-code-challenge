// Package api is responsible for defining the http endpoint handlers provided by the service.
package api

import (
	"net/http"

	"deus.ai-code-challenge/domain"
)

// Handlers returns all the service registered url and handler pairs
func Handlers(repo domain.VisitRepository) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /api/v1/unique-visitors":  buildUniqueVisitorForPageHandler(repo),
		"POST /api/v1/user-navigation": buildUserNavigationHandler(repo),
	}
}
