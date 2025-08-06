// Package repository is responsible for implementing an in-memory visit repository, optimized for the features requested in the code challenge.
package repository

import (
	"context"
	"deus.ai-code-challenge/domain"
)

const (
	kindStore = "store"
	kindCount = "count"
)

type visitorID = string

// InMemoryVisitRepository stores page visits in a structure optimised for the requirements provided
//   - data is a map of page urls (key) with their visitors (values)
//     *visitors is in itself a map of visitor id (key) with no values, (go doesn't provide a set data structure natively
//     but those can be mimicked by a map[KEY]struct{}. This ensures that visitors for a specific page are always unique
//   - count is a map of page urls (key) with the count of unique visitors (values)
//   - the counter is in itself a uint64 since a counter can never be negative and I'd expect a large number of unique visitor
//     *this lookup map ensures that reads are fast when handling a big number of visitors
//
// In terms of Big O notation this ensures both methods have an expected O(1) time complexity (exchanged for a higher space complexity)
//
// This version differs from the one provided in the main branch since it uses channels instead of mutexes to avoid data race problems.
// I'd advise against using this approach since a since mutex provides better performance and much cleaner code for this specific example
// (it's here just for me to have something more to show)
// Not that: a single channel is used to interact with the data (this ensures that interactions with the data are always ordered between different kinds)
type InMemoryVisitRepository struct {
	ch    chan query
	data  map[domain.PageURL]map[visitorID]struct{}
	count map[domain.PageURL]domain.Count
}

type query struct {
	kind  string
	store struct {
		ch    chan<- error
		visit domain.Visit
	}
	count struct {
		ch chan<- struct {
			count uint64
			err   error
		}
		pageURL domain.PageURL
	}
}

// NewVisitsInMemoryRepository is a constructor for the in-memory VisitsRepository
func NewVisitsInMemoryRepository(ctx context.Context) domain.VisitRepository {
	repo := &InMemoryVisitRepository{
		ch:    make(chan query),
		data:  make(map[domain.PageURL]map[visitorID]struct{}),
		count: make(map[domain.PageURL]domain.Count),
	}

	go repo.run(ctx)

	return repo
}

// Store is the public function that:
// - creates a `query` of type store
// - sends the query to the repo channel
// - waits for the response (within the channel sent) and returns it
func (i *InMemoryVisitRepository) Store(visit domain.Visit) error {
	ch := make(chan error)

	i.ch <- query{
		kind: kindStore,
		store: struct {
			ch    chan<- error
			visit domain.Visit
		}{ch: ch, visit: visit},
	}

	return <-ch
}

// Store ensures that unique visitor + page url are stored and accounted for when retrieving the counter for a page
func (i *InMemoryVisitRepository) store(visit domain.Visit) error {
	visitors, pageFound := i.data[visit.PageURL]
	if !pageFound {
		i.data[visit.PageURL] = map[visitorID]struct{}{
			visit.Visitor: {},
		}
		i.count[visit.PageURL] = 1

		return nil
	}

	_, visitorFound := visitors[visit.Visitor]
	if !visitorFound {
		i.data[visit.PageURL][visit.Visitor] = struct{}{}
		i.count[visit.PageURL]++
	}

	return nil
}

// CountUniqueVisitors simply reads the count map entry for the page url given
func (i *InMemoryVisitRepository) CountUniqueVisitors(url domain.PageURL) (domain.Count, error) {
	ch := make(chan struct {
		count uint64
		err   error
	})

	i.ch <- query{
		kind: kindCount,
		count: struct {
			ch chan<- struct {
				count domain.Count
				err   error
			}
			pageURL domain.PageURL
		}{ch: ch, pageURL: url},
	}

	resp := <-ch

	return resp.count, resp.err
}

// countUniqueVisitors simply reads the count map entry for the page url given
func (i *InMemoryVisitRepository) countUniqueVisitors(pageURL domain.PageURL) (domain.Count, error) {
	return i.count[pageURL], nil
}

// run iterates thought the incoming requests until the context signals that the process needs to close
func (i *InMemoryVisitRepository) run(ctx context.Context) {
	defer close(i.ch)

	for {
		select {
		case <-ctx.Done():
			return
		case q := <-i.ch:
			if q.kind == kindStore {
				err := i.store(q.store.visit)
				q.store.ch <- err
			}
			if q.kind == kindCount {
				count, err := i.countUniqueVisitors(q.count.pageURL)
				q.count.ch <- struct {
					count uint64
					err   error
				}{count: count, err: err}
			}
		}
	}
}
