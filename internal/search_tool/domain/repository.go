package domain

type SearchToolRepository interface {
	Search(filter SearchFilter) ([]SearchItem, error)
	Total(filter SearchFilter) (int64, error)
}
