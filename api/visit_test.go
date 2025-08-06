package api

import (
	"bytes"
	"deus.ai-code-challenge/domain"
	"deus.ai-code-challenge/service"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildUserNavigationHandler(t *testing.T) {
	type testCase struct {
		description        string
		input              string
		service            service.UserNavigationService
		expectedResponse   []byte
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description: "success",
			input:       `{"visitor_id": "id", "page_url": "url"}`,
			service: func(visit domain.Visit) error {
				return nil
			},
			expectedResponse:   []byte(``),
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "error: no visitor id provided",
			input:              `{"page_url": "url"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"missing request field: visitor id"}`),
		},
		{
			description:        "error: no page url provided",
			input:              `{"visitor_id": "id"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"missing request field: page url"}`),
		},
		{
			description:        "error: no body send",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"unable to read request body"}`),
		},
		{
			description: "error: call to repository fails",
			input:       `{"visitor_id": "id", "page_url": "url"}`,
			service: func(visit domain.Visit) error {
				return errors.New("failed to call repository")
			},
			expectedResponse:   []byte(`{"error":"failed to call repository"}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "error: page not found",
			input:       `{"visitor_id": "id", "page_url": "url"}`,
			service: func(visit domain.Visit) error {
				return service.ErrPageNotFound
			},
			expectedResponse:   []byte(`{"error":"page not found"}`),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "url", strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}

			h := buildUserNavigationHandler(tc.service)

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			resp := rr.Result()

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			body, _ := io.ReadAll(resp.Body)

			if !bytes.Equal(tc.expectedResponse, body) {
				t.Errorf("got %v, expected %v", string(body), string(tc.expectedResponse))
			}

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("got %v, expected %v", resp.StatusCode, tc.expectedStatusCode)
			}
		})
	}
}

func TestBuildUniqueVisitorForPageHandler(t *testing.T) {
	type testCase struct {
		description        string
		input              string
		service            service.UniqueVisitorForPageService
		expectedResponse   []byte
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description: "success",
			input:       `?pageUrl=url`,
			service: func(page domain.PageURL) (domain.Count, error) {
				return 10, nil
			},
			expectedResponse:   []byte(`{"unique_visitors":10}`),
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "error: no query param provided",
			input:              ``,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"missing query param: pageUrl"}`),
		},
		{
			description:        "error: no page url provided",
			input:              `?pageUrl=`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"missing query param: pageUrl"}`),
		},
		{
			description: "error: call to repository fails",
			input:       `?pageUrl=url`,
			service: func(page domain.PageURL) (domain.Count, error) {
				return 0, errors.New("failed to call repository")
			},
			expectedResponse:   []byte(`{"error":"failed to call repository"}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "error: page not found",
			input:       `?pageUrl=url`,
			service: func(page domain.PageURL) (domain.Count, error) {
				return 0, service.ErrPageNotFound
			},
			expectedResponse:   []byte(`{"error":"page not found"}`),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "url"+tc.input, nil)
			if err != nil {
				t.Fatal(err)
			}

			h := buildUniqueVisitorForPageHandler(tc.service)

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			resp := rr.Result()

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			body, _ := io.ReadAll(resp.Body)

			if !bytes.Equal(tc.expectedResponse, body) {
				t.Errorf("got %v, expected %v", string(body), string(tc.expectedResponse))
			}

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("got %v, expected %v", resp.StatusCode, tc.expectedStatusCode)
			}
		})
	}
}
