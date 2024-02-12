package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	Id        string     `json:"id,omitempty" gorm:"column:id;primaryKey"`
	UserName  *string    `json:"userName,omitempty" gorm:"column:user_name"`
	Email     *string    `json:"email,omitempty" gorm:"column:email"`
	BirthDate *time.Time `json:"birthDate,omitempty" gorm:"column:birth_date"`
	Phone     *string    `json:"phone,omitempty" gorm:"column:phone"`
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	CreatedBy *time.Time `json:"createdBy,omitempty" gorm:"column:created_by"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	UpdatedBy *time.Time `json:"updatedBy,omitempty" gorm:"column:updated_by"`
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
