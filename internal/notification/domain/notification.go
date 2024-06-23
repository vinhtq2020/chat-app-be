package domain

import "time"

type Notification struct {
	Id          string    `json:"id" gorm:"column:id"`
	RequestorId string    `json:"requestorId" gorm:"column:requestor_id"`
	Subscribers []string  `json:"subscribers" gorm:"column:subscribers"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Content     string    `json:"title" gorm:"column:title"`
	isRead      bool      `json:"isRead" gorm:"column:is_read"`
}
