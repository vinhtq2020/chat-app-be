package model

import "fmt"

type UserStatus string

var statusActive UserStatus = "A"
var statusUnActive UserStatus = "U"

func (u UserStatus) Value() string {
	return fmt.Sprintf("%v", u)
}

func StatusActive() *UserStatus {
	return &statusActive
}

func StatusUnActive() *UserStatus {
	return &statusUnActive
}
