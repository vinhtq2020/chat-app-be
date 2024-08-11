package domain

import (
	"fmt"
	"time"
)

type Friend struct {
	UserId1 string `json:"userId1" gorm:"column:user_id1"`
	UserId2 string `json:"userId2" gorm:"column:user_id2"`
	Status  string `json:"status" gorm:"column:status"`
}

type FriendRequest struct {
	Id             string              `json:"id" gorm:"column:id;primaryKey"`
	RequesterId    string              `json:"requesterId,omitempty" gorm:"column:requester_id"`
	RequesteeId    string              `json:"requesteeId,omitempty" gorm:"column:requestee_id" validate:"required"`
	Status         FriendRequestStatus `json:"status,omitempty" gorm:"column:status"`
	CreatedAt      time.Time           `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy      string              `json:"createdBy,omitempty" gorm:"column:created_by"`
	UpdatedAt      time.Time           `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	UpdatedBy      string              `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	NotificationId string              `json:"notificationId,omitempty" gorm:"column:notification_id"`
}

type FriendRequestStatus string

const AcceptAction = "accept"
const RejectAction = "reject"

const StatusAccept FriendRequestStatus = "A"
const StatusReject FriendRequestStatus = "R"
const StatusCancel FriendRequestStatus = "C"
const StatusPending FriendRequestStatus = "P"

func (u FriendRequestStatus) Value() string {
	return fmt.Sprintf("%v", u)
}
