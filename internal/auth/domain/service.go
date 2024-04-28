package domain

import (
	"context"
	"go-service/pkg/jwt"
	"go-service/pkg/validate"
)

type AuthService interface {
	Register(ctx context.Context, user Account) ([]validate.ErrorMsg, int64, error)
	LoginWithGoogle(ctx context.Context, email string) ([]validate.ErrorMsg, int64, *jwt.TokenData, error)
	Login(ctx context.Context, email string, password string, browser string, ipAdress string, deviceId string) ([]validate.ErrorMsg, *jwt.TokenData, error)
	Logout(ctx context.Context, userId string, browser string, ipAdress string, deviceId string) (int64, error)
	RefreshToken(ctx context.Context, userId string, browser string, ipAddress string, deviceId string) (int64, string, error)
}
