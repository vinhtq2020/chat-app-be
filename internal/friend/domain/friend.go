package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Relation struct {
	UserId1        string    `json:"userId1" gorm:"column:user_id1;primaryKey"`
	UserId2        string    `json:"userId2" gorm:"column:user_id2;primaryKey"`
	RelationType   string    `json:"relationType" gorm:"column:relation_type;primaryKey"`
	Status         string    `json:"status,omitempty" gorm:"column:status"`
	CreatedAt      time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy      string    `json:"createdBy,omitempty" gorm:"column:created_by"`
	UpdatedBy      string    `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	NotificationId string    `json:"notificationId,omitempty" gorm:"column:notification_id"`
}

func NewRelation(userId1 string, userId2 string, createdBy string, relationType RelationType, status Status, notificationId string) *Relation {
	now := time.Now()
	return &Relation{
		UserId1:        userId1,
		UserId2:        userId2,
		RelationType:   relationType.Value(),
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      createdBy,
		UpdatedBy:      createdBy,
		Status:         status.Value(),
		NotificationId: notificationId,
	}
}

type RelationType string

func (u RelationType) Value() string {
	return fmt.Sprintf("%v", u)
}

func (u *RelationType) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf(fmt.Sprintf("can't scan data to store in RelationType, expected []byte, got %T", v))
	}
	return json.Unmarshal(b, &u)
}

const NoneRelation RelationType = "none"
const FriendRelation RelationType = "friend"
const BlockRelation RelationType = "block"

type Status string

const AcceptAction = "accept"
const RejectAction = "reject"
const UnfriendAction = "unfriend"
const CancelAction = "cancel"

const StatusAccept Status = "A"
const StatusReject Status = "R"
const StatusCancel Status = "C"
const StatusPending Status = "P"
const StatusUnfriend Status = "U"

func (u Status) Value() string {
	return fmt.Sprintf("%v", u)
}

func (u *Status) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("can't scan data to store in Status")
	}
	return json.Unmarshal(b, &u)
}
