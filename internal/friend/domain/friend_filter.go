package domain

import "go-service/internal/utils/search"

type FriendFilter struct {
	search.SearchFilter
	UserId *string `json:"userId" gorm:"column:userId"`
}
