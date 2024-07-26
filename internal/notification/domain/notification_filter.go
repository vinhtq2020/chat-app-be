package domain

import (
	"time"
)

type NotificationFilter struct {
	SubscriberId *string
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	Visible      bool
}
