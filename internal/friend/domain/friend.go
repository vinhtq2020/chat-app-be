package domain

import (
	"fmt"
	"time"
)

type UserFriend struct {
	Id        string `json:"id" gorm:"column:id"`
	UserId    string `json:"userId1" gorm:"column:user_id"`
	FriendId  string `json:"userId2" gorm:"column:friend_id"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at"`
}

type FriendRequest struct {
	Id        string              `json:"id" gorm:"column:id"`
	Uid1      string              `json:"uid1" gorm:"column:uid1"`
	Uid2      string              `json:"uid2" gorm:"column:uid1"`
	Status    FriendRequestStatus `json:"status" gorm:"column:status"`
	CreatedAt time.Time           `json:"createdAt" gorm:"column:created_at"`
	CreatedBy string              `json:"createdBy" gorm:"column:created_by"`
	UpdatedAt time.Time           `json:"updatedAt" gorm:"column:updated_at"`
	UpdatedBy string              `json:"updatedBy" gorm:"column:updated_by"`
}

type FriendRequestStatus string

var StatusAccept FriendRequestStatus = "A"
var StatusReject FriendRequestStatus = "R"
var StatusCancel FriendRequestStatus = "C"
var StatusPending FriendRequestStatus = "P"

func (u FriendRequestStatus) Value() string {
	return fmt.Sprintf("%v", u)
}
