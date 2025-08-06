package service

import (
	"deus.ai-code-challenge/domain"
	"errors"
)

var (
	ErrPageNotFound = errors.New("page not found")
)

// buildUserNavigationService returns a function that verifies if the page visited is valid and if so, stores the visit
func buildUserNavigationService(visits domain.VisitsRepository, pages domain.PageRepository) func(visit domain.Visit) error {
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

// buildUniqueVisitorForPageService returns a function that verifies if the page visited is valid and if so, return the unique number of visitors for that page
func buildUniqueVisitorForPageService(visits domain.VisitsRepository, pages domain.PageRepository) func(page domain.PageURL) (domain.Count, error) {
	return func(page domain.PageURL) (domain.Count, error) {
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
