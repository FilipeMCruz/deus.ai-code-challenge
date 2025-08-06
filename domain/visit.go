// Package domain defines the "visit" concept and the interface required by the domain/features to manage "visits".
package domain

type Visit struct {
	Visitor string
	PageURL string
}

type PageURL = string
type Count = uint64

// VisitRepository is responsible for managing data related to user navigation according to the requirements provided
//
// Even though Store and CountUniqueVisitors can't fail when working with in-memory data structures, an error was added to the return
// so that we can better account for future changes (e.g. using redis instead of storing everything in memory so that there's no data lost when services are shutdown)
type VisitRepository interface {
	Store(visit Visit) error
	CountUniqueVisitors(url PageURL) (Count, error)
}
