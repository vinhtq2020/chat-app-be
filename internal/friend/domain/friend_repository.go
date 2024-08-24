package domain

import (
	"context"

	"gorm.io/gorm"
)

type FriendRepository interface {
	InTransaction(ctx context.Context, ex func(ctx context.Context, db *gorm.DB) (int64, error)) (int64, error)
	Load(ctx context.Context, userId1 string, userId2 string, relationType RelationType) (*Relation, error)
	Upsert(ctx context.Context, friend Relation) (int64, error)
	Patch(ctx context.Context, friendRelation map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
