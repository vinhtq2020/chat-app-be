package domain

import (
	"context"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	InTransaction(ctx context.Context, ex func(db *gorm.DB) (int64, error)) (int64, error)
	Load(ctx context.Context, browser string, ip string, deviceId string) (*RefreshToken, error)
	Insert(ctx context.Context, refreshToken RefreshToken) (int64, error)
	Patch(ctx context.Context, refreshToken map[string]interface{}) (int64, error)
	Delete(ctx context.Context, userId string, ipAddress string, deviceId string, browser string) (int64, error)
}
