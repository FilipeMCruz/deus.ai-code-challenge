package domain

type Visit struct {
	Visitor string
	PageURL string
}

// VisitsRepository is responsible for managing data related to user navigation according to the requirements provided
type VisitsRepository interface {
	Store(visit Visit) error
	CountUniqueVisitors(pageURL string) (int, error)
}
