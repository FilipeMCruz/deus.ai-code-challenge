// Package repository is responsible for implementing an:
// - in-memory visit repository, optimized for the features requested in the code challenge
// - in-memory page repository
package repository

import "deus.ai-code-challenge/domain"

type Repositories struct {
	Visits domain.VisitsRepository
	Pages  domain.PageRepository
}

// NewRepositories contains all repositories available, receives the list of valid pages
func NewRepositories(pages []domain.PageURL) Repositories {
	return Repositories{
		Visits: newVisitsInMemoryRepository(),
		Pages:  newPageInMemoryRepository(pages),
	}
}
