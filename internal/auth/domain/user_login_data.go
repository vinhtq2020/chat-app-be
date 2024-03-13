package domain

import (
	"go-service/pkg/model"
	"time"
)

type UserLoginData struct {
	Id           string            `json:"id,omitempty" gorm:"column:id"`
	Username     *string           `json:"username,omitempty" gorm:"column:user_name" validate:"required,min=6,max=50"`
	Email        *string           `json:"email,omitempty" gorm:"column:email" validate:"required,max=255"`
	Phone        *string           `json:"phone,omitempty" gorm:"column:phone" validate:"required"`
	Password     *string           `json:"password,omitempty" validate:"required_without=Provider"`
	PasswordHash *string           `json:"passwordHash,omitempty" gorm:"column:password_hash"`
	CreatedBy    *string           `json:"createdBy,omitempty" gorm:"column:created_by"`
	CreatedAt    *time.Time        `json:"createdAt,omitempty" gorm:"column:created_at"`
	UpdatedBy    *string           `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	UpdatedAt    *time.Time        `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	Status       *model.UserStatus `json:"status,omitempty" gorm:"column:status"`
	Provider     *model.Provider   `json:"provider,omitempty" gorm:"column:provider"`
	Version      int64             `json:"version,omitempty" gorm:"column:version"`
}

type ExistingField struct {
	ExistingField string `json:"existingField" gorm:"column:existing_field"`
	ExistingValue string `json:"existingValue" gorm:"column:existing_value"`
}
