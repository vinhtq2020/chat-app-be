package domain

import (
	"time"
)

type Notification struct {
	Id          string       `json:"id" gorm:"column:id;primaryKey"`
	Requester   Requester    `json:"requester,omitempty" gorm:"column:requester"`
	Title       string       `json:"title,omitempty" gorm:"column:title"`
	Content     string       `json:"content,omitempty" gorm:"column:content"`
	Subscribers []Subscriber `json:"subscribers,omitempty" gorm:"column:subscribers"`
	Type        string       `json:"type,omitempty" gorm:"column:type"`
	URL         *string      `json:"url,omitempty" gorm:"column:url"`
	CreatedAt   time.Time    `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy   string       `json:"createdBy,omitempty" gorm:"column:created_by"`
	UpdatedAt   time.Time    `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	UpdatedBy   string       `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	Visible     bool         `json:"visible,omitempty" gorm:"column:visible"`
}

type Subscriber struct {
	Id     string `json:"id" gorm:"column:subscriber_id"`
	IsRead bool   `json:"isRead,omitempty" gorm:"column:is_read"`
}

type Requester struct {
	Id        string  `json:"id" gorm:"column:id"`
	Name      string  `json:"name,omitempty" gorm:"column:name"`
	AvatarURL *string `json:"avatarURL,omitempty" gorm:"column:avatar_url"`
}

const NotificationInvite = "invite"
const NotificationInform = "inform"
