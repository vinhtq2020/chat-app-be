package domain

import (
	"encoding/json"
	"errors"
	"fmt"
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
}

type Subscriber struct {
	Id     string `json:"id"`
	Readed bool   `json:"readed"`
}

func (u *Subscriber) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("can't scan data to store in Subscriber")
	}
	return json.Unmarshal(b, &u)
}

func (u Subscriber) Value() string {
	return fmt.Sprintf("%v", u)
}

type Requester struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatarURL"`
}

func NewRequester(Id string, name string, avatarURL *string) Requester {
	return Requester{
		Id:        Id,
		Name:      name,
		AvatarURL: avatarURL,
	}
}

func (u *Requester) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("can't scan data to store in Requester")
	}
	return json.Unmarshal(b, &u)
}

func (u Requester) Value() string {
	return fmt.Sprintf("%v", u)
}

const NotificationAddFriend = "addfriend"
const NotificationInform = "inform"
