package api

import (
	"deus.ai-code-challenge/domain"
	"deus.ai-code-challenge/service"
	"encoding/json"
	"io"
	"net/http"
)

// buildUserNavigationHandler provides an http handler responsible for storing a new visit
func buildUserNavigationHandler(repository domain.VisitsRepository, pages domain.PageRepository) http.HandlerFunc {
	serv := service.BuildUserNavigationService(repository, pages)

	type requestBody struct {
		VisitorId string `json:"visitor_id,omitempty"`
		PageURL   string `json:"page_url,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		i := &requestBody{}

		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			writeError(w, errUnmarshallRequest, http.StatusBadRequest)

			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		if i.VisitorId == "" {
			writeError(w, errMissingFieldPrefix+"visitor id", http.StatusBadRequest)

			return
		}

		if i.PageURL == "" {
			writeError(w, errMissingFieldPrefix+"page url", http.StatusBadRequest)

			return
		}

		err = serv(domain.Visit{
			Visitor: i.VisitorId,
			PageURL: i.PageURL,
		})
		if err != nil {
			writeError(w, err.Error(), getStatusCode(err))

			return
		}
	}
}

// buildUniqueVisitorForPageHandler provides an http.Handler responsible for providing the unique number of visitor
// for a specific page
func buildUniqueVisitorForPageHandler(repository domain.VisitsRepository, pages domain.PageRepository) http.HandlerFunc {
	serv := service.BuildUniqueVisitorForPageService(repository, pages)

	queryParamKey := "pageUrl"

	type responseBody struct {
		UniqueVisitors uint64 `json:"unique_visitors"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		pageURL := r.URL.Query().Get(queryParamKey)
		if pageURL == "" {
			writeError(w, errMissingParamPrefix+queryParamKey, http.StatusBadRequest)

			return
		}

		numberOfUniqueVisitors, err := serv(domain.PageURL(pageURL))
		if err != nil {
			writeError(w, err.Error(), getStatusCode(err))

			return
		}

		b, err := json.Marshal(responseBody{UniqueVisitors: numberOfUniqueVisitors})
		if err != nil {
			writeError(w, errMarshallResponse, http.StatusInternalServerError)

			return
		}

		_, _ = w.Write(b)
	}
}
