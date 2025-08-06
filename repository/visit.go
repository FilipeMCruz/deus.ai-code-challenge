// Package repository is responsible for implementing an in-memory visit repository, optimized for the features requested in the code challenge.
package repository

import (
	"deus.ai-code-challenge/domain"
	"sync"
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
type InMemoryVisitRepository struct {
	m     sync.RWMutex
	data  map[domain.PageURL]map[visitorID]struct{}
	count map[domain.PageURL]domain.Count
}

// NewVisitsInMemoryRepository is a constructor for the in-memory VisitRepository
func NewVisitsInMemoryRepository() domain.VisitRepository {
	return &InMemoryVisitRepository{
		data:  make(map[domain.PageURL]map[visitorID]struct{}),
		count: make(map[domain.PageURL]domain.Count),
	}
}

// Store ensures that unique visitor + page url are stored and accounted for when retrieving the counter for a page
func (i *InMemoryVisitRepository) Store(visit domain.Visit) error {
	i.m.Lock()
	defer i.m.Unlock()

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
	i.m.RLock()
	defer i.m.RUnlock()

	return i.count[url], nil
}
