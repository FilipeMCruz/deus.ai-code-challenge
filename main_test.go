package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	type req struct {
		method       string
		url          string
		body         string
		expectedBody string
		expectedCode int
	}

	type testCase struct {
		description string
		reqs        []req
		err         error
	}

	testCases := []testCase{
		{
			description: "unique-visitors -> user-navigation -> unique-visitors -> user-navigation -> unique-visitors -> user-navigation -> unique-visitors",
			reqs: []req{
				{
					method:       http.MethodGet,
					url:          "/api/v1/unique-visitors?pageUrl=url",
					expectedCode: http.StatusOK,
					expectedBody: `{"unique_visitors":0}`,
				},
				{
					method:       http.MethodPost,
					url:          "/api/v1/user-navigation",
					body:         `{"visitor_id": "id", "page_url": "url"}`,
					expectedCode: http.StatusOK,
				},
				{
					method:       http.MethodGet,
					url:          "/api/v1/unique-visitors?pageUrl=url",
					expectedCode: http.StatusOK,
					expectedBody: `{"unique_visitors":1}`,
				},
				{
					method:       http.MethodPost,
					url:          "/api/v1/user-navigation",
					body:         `{"visitor_id": "id2", "page_url": "url"}`,
					expectedCode: http.StatusOK,
				},
				{
					method:       http.MethodPost,
					url:          "/api/v1/user-navigation",
					body:         `{"visitor_id": "id2", "page_url": "url"}`,
					expectedCode: http.StatusOK,
				},
				{
					method:       http.MethodPost,
					url:          "/api/v1/user-navigation",
					body:         `{"visitor_id": "id3", "page_url": "url"}`,
					expectedCode: http.StatusOK,
				},
				{
					method:       http.MethodGet,
					url:          "/api/v1/unique-visitors?pageUrl=url",
					expectedCode: http.StatusOK,
					expectedBody: `{"unique_visitors":3}`,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx, stop := context.WithCancel(context.Background())
			defer stop()

			port, err := GetFreePort()
			if err != nil {
				t.Fatal(err)
			}

			go func() {
				err := start(ctx, stop, port)

				if !reflect.DeepEqual(tc.err, err) {
					t.Errorf("got %v, expected %v", err, tc.err)
				}
			}()

			time.Sleep(time.Second)

			for _, req := range tc.reqs {
				r, _ := http.NewRequest(req.method, "http://localhost:"+strconv.Itoa(port)+req.url, strings.NewReader(req.body))

				resp, err := http.DefaultClient.Do(r)
				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != req.expectedCode {
					t.Errorf("got %d, expected %d", resp.StatusCode, req.expectedCode)
				}

				body, _ := io.ReadAll(resp.Body)
				if string(body) != req.expectedBody {
					t.Errorf("got %s, expected %s", string(body), req.expectedBody)
				}
			}
		})
	}
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err == nil {
		l, err := net.ListenTCP("tcp", addr)
		if err == nil {
			defer func(l *net.TCPListener) {
				_ = l.Close()
			}(l)
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}

	return -1, err
}
