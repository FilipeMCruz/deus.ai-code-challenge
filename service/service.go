// Package service is responsible for describing and executing the business rules and processes. It should only be dependent on the domain
package service

import (
	"deus.ai-code-challenge/domain"
	"deus.ai-code-challenge/repository"
)

type UserNavigationService func(visit domain.Visit) error
type UniqueVisitorForPageService func(page domain.PageURL) (domain.Count, error)

type Services struct {
	UserNavigationService       UserNavigationService
	UniqueVisitorForPageService UniqueVisitorForPageService
}

// NewServices contains all services available, receives all repositories available
func NewServices(repositories repository.Repositories) Services {
	return Services{
		UserNavigationService:       buildUserNavigationService(repositories.Visits, repositories.Pages),
		UniqueVisitorForPageService: buildUniqueVisitorForPageService(repositories.Visits, repositories.Pages),
	}
}
