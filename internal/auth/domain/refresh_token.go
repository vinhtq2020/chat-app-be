package domain

import "time"

type RefreshToken struct {
	UserId    string        `json:"userId" gorm:"column:user_id"`
	DeviceId  string        `json:"deviceId" gorm:"column:device_id"`
	Browser   string        `json:"browser" gorm:"column:browser"`
	IPAddress string        `json:"ip_address" gorm:"column:ip_address"`
	Token     string        `json:"token" gorm:"column:token"`
	Expiry    time.Duration `json:"expiry" gorm:"column:expiry"`
	CreatedAt time.Time     `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time     `json:"updatedAt" gorm:"column:updated_at"`
}
