// Package service is responsible for describing and executing the business rules and processes. It should only be dependent on the domain
package service

import (
	"deus.ai-code-challenge/domain"
	"errors"
)

var (
	ErrPageNotFound = errors.New("page not found")
)

// BuildUserNavigationService returns a function that verifies if the page visited is valid and if so, stores the visit
func BuildUserNavigationService(visits domain.VisitsRepository, pages domain.PageRepository) func(visit domain.Visit) error {
	return func(visit domain.Visit) error {
		found, err := pages.Exists(domain.PageURL(visit.PageURL))
		if err != nil {
			return err
		}

		if !found {
			return ErrPageNotFound
		}

		return visits.Store(visit)
	}
}

// BuildUniqueVisitorForPageService returns a function that verifies if the page visited is valid and if so, return the unique number of visitors for that page
func BuildUniqueVisitorForPageService(visits domain.VisitsRepository, pages domain.PageRepository) func(page domain.PageURL) (uint64, error) {
	return func(page domain.PageURL) (uint64, error) {
		found, err := pages.Exists(page)
		if err != nil {
			return 0, err
		}

		if !found {
			return 0, ErrPageNotFound
		}

		return visits.CountUniqueVisitors(page)
	}
}
