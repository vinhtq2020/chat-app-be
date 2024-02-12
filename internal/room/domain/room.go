package domain

import (
	"go-service/internal/user/domain"
	"time"
)

type Room struct {
	Id        string        `json:"id,omitempty" gorm:"column:id;primaryKey"`
	Version   *int32        `json:"version,omitempty" gorm:"column:version"`
	Name      *string       `json:"name,omitempty" gorm:"column:name"`
	Members   []domain.User `json:"members,omitempty" gorm:"column:members"`
	UpdatedAt *time.Time    `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	UpdatedBy *string       `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	CreatedAt *time.Time    `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy *string       `json:"createdBy,omitempty" gorm:"column:created_by"`
}
