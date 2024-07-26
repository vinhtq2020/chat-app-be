package domain

import (
	"fmt"
	"time"
)

type FriendRequest struct {
	Id          string              `json:"id" gorm:"column:id"`
	RequesterId string              `json:"requesterId,omitempty" gorm:"column:requester_id"`
	RequesteeId string              `json:"requesteeId,omitempty" gorm:"column:requestee_id"`
	Status      FriendRequestStatus `json:"status,omitempty" gorm:"column:status"`
	CreatedAt   time.Time           `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy   string              `json:"createdBy,omitempty" gorm:"column:created_by"`
	UpdatedAt   time.Time           `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	UpdatedBy   string              `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	Action      string              `json:"action,omitempty"`
}

type FriendRequestStatus string

var StatusAccept FriendRequestStatus = "A"
var StatusReject FriendRequestStatus = "R"
var StatusCancel FriendRequestStatus = "C"
var StatusPending FriendRequestStatus = "P"

func (u FriendRequestStatus) Value() string {
	return fmt.Sprintf("%v", u)
}
