package service

import (
	"deus.ai-code-challenge/domain"
	"errors"
	"reflect"
	"testing"
)

func TestBuildUserNavigationHandler(t *testing.T) {
	type testCase struct {
		description    string
		input          domain.Visit
		mockRepoFunc   func(visit domain.Visit) error
		existsRepoFunc func(url domain.PageURL) (bool, error)
		output         error
	}

	testCases := []testCase{
		{
			description: "success",
			input:       domain.Visit{PageURL: "url", Visitor: "id"},
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
			output: nil,
		},
		{
			description: "error: call to repository fails",
			input:       domain.Visit{PageURL: "url", Visitor: "id"},
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
			output: errors.New("failed to call repository"),
		},
		{
			description: "error: page not found",
			input:       domain.Visit{PageURL: "url", Visitor: "id"},
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return false, nil
			},
			output: ErrPageNotFound,
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

			serv := buildUserNavigationService(mockRepo, mockPageRepo)

			err := serv(tc.input)

			if !reflect.DeepEqual(err, tc.output) {
				t.Errorf("err = %#v, want %#v", err, tc.output)
			}
		})
	}
}

func TestBuildUniqueVisitorForPageHandler(t *testing.T) {
	type testCase struct {
		description    string
		input          domain.PageURL
		mockRepoFunc   func(pageURL domain.PageURL) (domain.Count, error)
		existsRepoFunc func(url domain.PageURL) (bool, error)
		output         domain.Count
		outputErr      error
	}

	testCases := []testCase{
		{
			description: "success",
			input:       domain.PageURL("url"),
			mockRepoFunc: func(pageURL domain.PageURL) (domain.Count, error) {
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
			output: 10,
		},
		{
			description: "error: call to repository fails",
			input:       domain.PageURL("url"),
			mockRepoFunc: func(pageURL domain.PageURL) (domain.Count, error) {
				if pageURL != "url" {
					t.Errorf("pageURL = %v, want %v", pageURL, "url")
				}

				return 0, err
			},
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return true, nil
			},
			outputErr: err,
		},
		{
			description: "error: page not found",
			input:       domain.PageURL("url"),
			existsRepoFunc: func(url domain.PageURL) (bool, error) {
				if url != "url" {
					t.Errorf("url = %v, want %v", url, "url")
				}

				return false, nil
			},
			outputErr: ErrPageNotFound,
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

			serv := buildUniqueVisitorForPageService(mockRepo, mockPageRepo)

			output, outputErr := serv(tc.input)

			if !errors.Is(outputErr, tc.outputErr) {
				t.Errorf("output = %#v, want %#v", outputErr, tc.outputErr)
			}

			if !reflect.DeepEqual(tc.output, output) {
				t.Errorf("output = %#v, want %#v", output, tc.output)
			}
		})
	}
}
