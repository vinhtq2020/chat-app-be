package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	Id         string     `json:"id,omitempty" gorm:"column:id;primaryKey"`
	UserName   *string    `json:"userName,omitempty" gorm:"column:user_name"`
	FirstName  *string    `json:"email,omitempty" gorm:"column:first_name"`
	LastName   *string    `json:"lastName,omitempty" gorm:"column:last_name"`
	MiddleName *string    `json:"middleName,omitempty" gorm:"column:middle_name"`
	BirthDate  *time.Time `json:"birthDate,omitempty" gorm:"column:birth_date"`
	CreatedAt  *time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy  *time.Time `json:"createdBy,omitempty" gorm:"column:created_by"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	UpdatedBy  *time.Time `json:"updatedBy,omitempty" gorm:"column:updated_by"`
	Version    int64      `json:"version,omitempty" gorm:"column:version"`
}

func (u *User) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("can't scan data to store in User")
	}
	return json.Unmarshal(b, &u)
}

func (u User) Value() (driver.Value, error) {
	return json.Marshal(u)
}
