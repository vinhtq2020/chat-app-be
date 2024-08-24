package domain

import (
	"fmt"
	"time"
)

type Relation struct {
	UserId1      string       `json:"userId1" gorm:"column:user_id1;primaryKey"`
	UserId2      string       `json:"userId2" gorm:"column:user_id2;primaryKey"`
	RelationType RelationType `json:"relationType" gorm:"column:relation_type;primaryKey"`
	Status       Status       `json:"status" gorm:"column:status"`
	CreatedAt    time.Time    `json:"createdAt" gorm:"column:created_at"`
	CreatedBy    string       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy    string       `json:"updatedBy" gorm:"column:updated_by"`
	UpdatedAt    time.Time    `json:"updatedAt" gorm:"column:updated_at"`
}

func NewRelation(userId1 string, userId2 string, createdBy string, relationType RelationType, status Status) Relation {
	now := time.Now()
	return Relation{
		UserId1:      userId1,
		UserId2:      userId2,
		RelationType: relationType,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
		Status:       status,
	}
}

type RelationType string

func (u RelationType) Value() string {
	return fmt.Sprintf("%v", u)
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
