package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"deus.ai-code-challenge/domain"
)

// buildUserNavigationHandler provides an http handler responsible for storing a new visit
func buildUserNavigationHandler(repository domain.VisitRepository) http.HandlerFunc {
	type requestBody struct {
		VisitorId string `json:"visitor_id,omitempty"`
		PageURL   string `json:"page_url,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		i := &requestBody{}

		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			writeError(w, newErrUnmarshallRequest())

			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		if i.VisitorId == "" {
			writeError(w, newErrMissingFieldPrefix("visitor id"))

			return
		}

		if i.PageURL == "" {
			writeError(w, newErrMissingFieldPrefix("page url"))

			return
		}

		_, err = url.Parse(i.PageURL)
		if err != nil {
			writeError(w, newErrInvalidPageURL(i.PageURL))

			return
		}

		err = repository.Store(domain.Visit{
			Visitor: i.VisitorId,
			PageURL: i.PageURL,
		})
		if err != nil {
			writeError(w, err)

			return
		}
	}
}

// buildUniqueVisitorForPageHandler provides an http.Handler responsible for providing the unique number of visitor
// for a specific page
func buildUniqueVisitorForPageHandler(repository domain.VisitRepository) http.HandlerFunc {
	queryParamKey := "pageUrl"

	type responseBody struct {
		UniqueVisitors uint64 `json:"unique_visitors"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		pageURL := r.URL.Query().Get(queryParamKey)
		if pageURL == "" {
			writeError(w, newErrMissingParamPrefix(queryParamKey))

			return
		}

		_, err := url.Parse(pageURL)
		if err != nil {
			writeError(w, newErrInvalidPageURL(pageURL))

			return
		}

		numberOfUniqueVisitors, err := repository.CountUniqueVisitors(pageURL)
		if err != nil {
			writeError(w, err)

			return
		}

		b, err := json.Marshal(responseBody{UniqueVisitors: numberOfUniqueVisitors})
		if err != nil {
			writeError(w, newErrMarshallResponse())

			return
		}

		_, _ = w.Write(b)
	}
}
