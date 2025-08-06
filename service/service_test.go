package service

import (
	"deus.ai-code-challenge/domain"
	"errors"
	"testing"
)

var err = errors.New("failed to call repository")

type mockPageRepository struct {
	t          *testing.T
	existsFunc func(url domain.PageURL) (bool, error)
}

func (m *mockPageRepository) Exists(url domain.PageURL) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(url)
	}

	m.t.Fatal("mockVisitRepository storeFunc is nil")
	return false, nil
}

type mockVisitRepository struct {
	t                   *testing.T
	storeFunc           func(domain.Visit) error
	countUniqueVisitors func(pageURL domain.PageURL) (domain.Count, error)
}

func (m *mockVisitRepository) Store(visit domain.Visit) error {
	if m.storeFunc != nil {
		return m.storeFunc(visit)
	}

	m.t.Fatal("mockVisitRepository storeFunc is nil")
	return nil
}

func (m *mockVisitRepository) CountUniqueVisitors(pageURL domain.PageURL) (domain.Count, error) {
	if m.countUniqueVisitors != nil {
		return m.countUniqueVisitors(pageURL)
	}

	m.t.Fatal("mockVisitRepository CountUniqueVisitors is nil")
	return 0, nil
}
