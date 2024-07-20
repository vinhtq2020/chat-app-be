package friend_domain

import (
	"fmt"
	"time"
)

type FriendRequest struct {
	Id          string               `json:"id" gorm:"column:id"`
	RequesterId string               `json:"requesterId" gorm:"column:requester_id"`
	RequesteeId string               `json:"requesteeId" gorm:"column:requestee_id"`
	Status      *FriendRequestStatus `json:"status" gorm:"column:status"`
	CreatedAt   time.Time            `json:"createdAt" gorm:"column:created_at"`
	CreatedBy   string               `json:"createdBy" gorm:"column:created_by"`
	UpdatedAt   time.Time            `json:"updatedAt" gorm:"column:updated_at"`
	UpdatedBy   string               `json:"updatedBy" gorm:"column:updated_by"`
}

type FriendRequestStatus string

var StatusAccept FriendRequestStatus = "A"
var StatusReject FriendRequestStatus = "R"
var StatusCancel FriendRequestStatus = "C"
var StatusPending FriendRequestStatus = "P"

func (u FriendRequestStatus) Value() string {
	return fmt.Sprintf("%v", u)
}
