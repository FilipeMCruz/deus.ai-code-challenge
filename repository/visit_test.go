package repository

import (
	"deus.ai-code-challenge/domain"
	"sync"
	"testing"
)

func TestInMemoryRepository(t *testing.T) {
	type count struct {
		pageURL       domain.PageURL
		expectedCount uint64
	}

	type input struct {
		store domain.Visit
		count count
	}

	type testCase struct {
		description string
		inputs      []input
	}

	testCases := []testCase{
		{
			description: "basic insert & search",
			inputs: []input{
				{
					store: domain.Visit{
						Visitor: "id",
						PageURL: "url",
					},
				},
				{
					count: count{
						pageURL:       "url",
						expectedCount: 1,
					},
				},
			},
		},
		{
			description: "multiple inserts for same visitor+page & search",
			inputs: []input{
				{
					store: domain.Visit{
						Visitor: "id",
						PageURL: "url",
					},
				},
				{
					store: domain.Visit{
						Visitor: "id",
						PageURL: "url",
					},
				},
				{
					count: count{
						pageURL:       "url",
						expectedCount: 1,
					},
				},
			},
		},
		{
			description: "multiple inserts & search",
			inputs: []input{
				{
					store: domain.Visit{
						Visitor: "id",
						PageURL: "url",
					},
				},
				{
					store: domain.Visit{
						Visitor: "id2",
						PageURL: "url",
					},
				},
				{
					count: count{
						pageURL:       "url",
						expectedCount: 2,
					},
				},
			},
		},
		{
			description: "search for un-visited page",
			inputs: []input{
				{
					count: count{
						pageURL:       "url",
						expectedCount: 0,
					},
				},
			},
		},
		{
			description: "multiple inserts & searches",
			inputs: []input{
				{
					store: domain.Visit{
						Visitor: "id",
						PageURL: "url",
					},
				},
				{
					store: domain.Visit{
						Visitor: "id2",
						PageURL: "url2",
					},
				},
				{
					store: domain.Visit{
						Visitor: "id2",
						PageURL: "url",
					},
				},
				{
					count: count{
						pageURL:       "url",
						expectedCount: 2,
					},
				},
				{
					count: count{
						pageURL:       "url2",
						expectedCount: 1,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r := NewVisitsInMemoryRepository()

			for _, input := range tc.inputs {
				if input.store.PageURL == "" {
					counter, err := r.CountUniqueVisitors(input.count.pageURL)
					if err != nil {
						t.Fatal("unexpected error", err)
					}

					if counter != input.count.expectedCount {
						t.Errorf("got %v, expected %v", counter, input.count.expectedCount)
					}
				} else {
					err := r.Store(input.store)
					if err != nil {
						t.Fatal("unexpected error", err)
					}
				}
			}
		})
	}
}

func TestInMemoryRepositoryConcurrency(t *testing.T) {
	type testCase struct {
		description    string
		inputs         []domain.Visit
		expectedCounts map[domain.PageURL]domain.Count
	}

	testCases := []testCase{
		{
			description: "multiple inserts",
			inputs: []domain.Visit{
				{
					Visitor: "id1",
					PageURL: "url",
				},
				{
					Visitor: "id2",
					PageURL: "url",
				},
				{
					Visitor: "id3",
					PageURL: "url",
				},
				{
					Visitor: "id4",
					PageURL: "url",
				},
				{
					Visitor: "id5",
					PageURL: "url",
				},
			},
			expectedCounts: map[domain.PageURL]domain.Count{
				"url": 5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r := NewVisitsInMemoryRepository()

			var wg = sync.WaitGroup{}
			wg.Add(len(tc.inputs))
			for _, i := range tc.inputs {
				go func() {
					err := r.Store(i)
					if err != nil {
						t.Error("unexpected error", err)
					}
					wg.Done()
				}()
			}

			wg.Wait()
			for k, v := range tc.expectedCounts {
				counter, err := r.CountUniqueVisitors(k)
				if err != nil {
					t.Error("unexpected error", err)
				}

				if counter != v {
					t.Errorf("got %v, expected %v", counter, v)
				}
			}
		})
	}
}
