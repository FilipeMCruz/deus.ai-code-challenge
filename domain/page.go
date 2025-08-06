package domain

type PageURL string

// PageRepository is responsible for managing data related to available page within the website
type PageRepository interface {
	Exists(pageURL PageURL) (bool, error)
}
