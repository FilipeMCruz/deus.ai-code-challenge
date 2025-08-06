package api

import (
	"bytes"
	"deus.ai-code-challenge/domain"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockPageRepository struct {
	t          *testing.T
	existsFunc func(url domain.PageURL) (bool, error)
}

func (m *mockPageRepository) Exists(url domain.PageURL) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(url)
	}

	m.t.Fatal("mockVisitRepository storeFunc is nil")
	return false, nil
}

type mockVisitRepository struct {
	t                   *testing.T
	storeFunc           func(domain.Visit) error
	countUniqueVisitors func(pageURL domain.PageURL) (uint64, error)
}

func (m *mockVisitRepository) Store(visit domain.Visit) error {
	if m.storeFunc != nil {
		return m.storeFunc(visit)
	}

	m.t.Fatal("mockVisitRepository storeFunc is nil")
	return nil
}

func (m *mockVisitRepository) CountUniqueVisitors(pageURL domain.PageURL) (uint64, error) {
	if m.countUniqueVisitors != nil {
		return m.countUniqueVisitors(pageURL)
	}

	m.t.Fatal("mockVisitRepository CountUniqueVisitors is nil")
	return 0, nil
}

func TestBuildUserNavigationHandler(t *testing.T) {
	type testCase struct {
		description        string
		input              string
		mockRepoFunc       func(visit domain.Visit) error
		existsRepoFunc     func(url domain.PageURL) (bool, error)
		expectedResponse   []byte
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description: "success",
			input:       `{"visitor_id": "id", "page_url": "url"}`,
			mockRepoFunc: func(visit domain.Visit) error {
				if visit.PageURL != "url" {
					t.Errorf("visit.PageURL = %v, want %v", visit.PageURL, "url")
				}
				if visit.Visitor != "id" {
					t.Errorf("visit.Visitor = %v, want %v", visit.Visitor, "id")
				}

				return nil
			},
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return true, nil
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
			mockRepoFunc: func(visit domain.Visit) error {
				if visit.PageURL != "url" {
					t.Errorf("visit.PageURL = %v, want %v", visit.PageURL, "url")
				}
				if visit.Visitor != "id" {
					t.Errorf("visit.Visitor = %v, want %v", visit.Visitor, "id")
				}

				return errors.New("failed to call repository")
			},
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return true, nil
			},
			expectedResponse:   []byte(`{"error":"failed to call repository"}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "error: page not found",
			input:       `{"visitor_id": "id", "page_url": "url"}`,
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return false, nil
			},
			expectedResponse:   []byte(`{"error":"page not found"}`),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &mockVisitRepository{
				t:         t,
				storeFunc: tc.mockRepoFunc,
			}

			mockPageRepo := &mockPageRepository{
				t:          t,
				existsFunc: tc.existsRepoFunc,
			}

			req, err := http.NewRequest(http.MethodPost, "url", strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}

			h := buildUserNavigationHandler(mockRepo, mockPageRepo)

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
		mockRepoFunc       func(pageURL domain.PageURL) (uint64, error)
		existsRepoFunc     func(url domain.PageURL) (bool, error)
		expectedResponse   []byte
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description: "success",
			input:       `?pageUrl=url`,
			mockRepoFunc: func(pageURL domain.PageURL) (uint64, error) {
				if pageURL != "url" {
					t.Errorf("pageURL = %v, want %v", pageURL, "url")
				}

				return 10, nil
			},
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return true, nil
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
			mockRepoFunc: func(pageURL domain.PageURL) (uint64, error) {
				if pageURL != "url" {
					t.Errorf("pageURL = %v, want %v", pageURL, "url")
				}

				return 0, errors.New("failed to call repository")
			},
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return true, nil
			},
			expectedResponse:   []byte(`{"error":"failed to call repository"}`),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "error: page not found",
			input:       `?pageUrl=url`,
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return false, nil
			},
			expectedResponse:   []byte(`{"error":"page not found"}`),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := &mockVisitRepository{
				t:                   t,
				countUniqueVisitors: tc.mockRepoFunc,
			}

			mockPageRepo := &mockPageRepository{
				t:          t,
				existsFunc: tc.existsRepoFunc,
			}

			req, err := http.NewRequest(http.MethodPost, "url"+tc.input, nil)
			if err != nil {
				t.Fatal(err)
			}

			h := buildUniqueVisitorForPageHandler(mockRepo, mockPageRepo)

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
