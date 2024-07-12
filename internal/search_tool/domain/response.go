package domain

type SearchResponse struct {
	List  []SearchItem `json:"list"`
	Total int64        `json:"total"`
}
