package search

type SearchFilter struct {
	Q     *string  `json:"q,omitempty"`
	Page  *int64   `json:"page,omitempty"`
	Limit *int64   `json:"limit,omitempty"`
	Sorts []string `json:"sorts,omitempty"`
}
