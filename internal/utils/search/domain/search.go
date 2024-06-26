package domain

type SearchFilter struct {
	Q     *string  `json:"q"`
	Page  *int64   `json:"page"`
	Limit *int64   `json:"limit"`
	Sorts []string `json:"sorts"`
}
