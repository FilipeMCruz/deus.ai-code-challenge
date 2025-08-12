package infrastructure

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWrappers(t *testing.T) {
	type testCase struct {
		description        string
		handler            http.Handler
		input              *http.Request
		expectedResponse   []byte
		expectedStatusCode int
		expectedHeaders    http.Header
	}

	testCases := []testCase{
		{
			description: "recovery from panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				panic(nil)
			}),
			input: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "/", nil)
				return r
			}(),
			expectedResponse:   nil,
			expectedStatusCode: http.StatusInternalServerError,
			expectedHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		},
		{
			description: "basic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				_, _ = w.Write([]byte(`{"key":"value"}`))
			}),
			input: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "/", nil)
				return r
			}(),
			expectedResponse:   []byte(`{"key":"value"}`),
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			handler := Wrap(tc.handler)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tc.input)
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

			if !reflect.DeepEqual(tc.expectedHeaders, resp.Header) {
				t.Errorf("got %v, expected %v", resp.Header, tc.expectedHeaders)
			}
		})
	}
}
