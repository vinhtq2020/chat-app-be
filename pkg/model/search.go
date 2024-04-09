package model

type SearchResult struct {
	List  interface{} `json:"list"`
	Total interface{} `json:"total"`
}
