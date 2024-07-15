package domain

type SearchResult interface{}

type SearchItem struct {
	Id           string  `json:"id" gorm:"column:id"`
	Name         *string `json:"name,omitempty" gorm:"column:user_name"`
	Avatar       *string `json:"avatar,omitempty" gorm:"column:avatar_url"`
	Carrier      *string `json:"carrier,omitempty" gorm:"column:carrier"`
	FriendStatus *string `json:"friendStatus,omitempty" gorm:"column:friend_status"`
}
