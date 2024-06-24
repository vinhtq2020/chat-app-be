package domain

import "time"

type Notification struct {
	Id          string       `json:"id" gorm:"column:id"`
	RequestorId string       `json:"requestorId" gorm:"column:requestor_id"`
	Subscribers []Subscriber `json:"subscribers" gorm:"column:subscribers"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"column:updated_at"`
	Content     string       `json:"title" gorm:"column:content"`
}

type Subscriber struct {
	SubscriberId string `json:"subscriberId" gorm:"column:subscriber_id"`
	IsRead       bool   `json:"isRead" gorm:"column:is_read"`
}
