package domain

import (
	"go-service/pkg/model"
	"time"
)

type UserLoginData struct {
	Id           string            `json:"id,omitempty" gorm:"column:id"`
	Username     *string           `json:"username,omitempty" gorm:"column:username" validate:"required"`
	Email        *string           `json:"email,omitempty" gorm:"column:email" validate:"required"`
	Password     *string           `json:"password,omitempty" validate:"required_without=Provider"`
	PasswordHash []byte            `json:"passwordHash,omitempty" gorm:"column:password_hash"`
	CreatedBy    *string           `json:"createdBy,omitempty" gorm:"column:created_by"`
	CreatedAt    *time.Time        `json:"createdAt,omitempty" gorm:"column:created_at"`
	UpdatedBy    *string           `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	UpdatedAt    *time.Time        `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	Status       *model.UserStatus `json:"status,omitempty" gorm:"column:status"`
	Provider     *model.Provider   `json:"provider,omitempty" gorm:"column:provider"`
}
