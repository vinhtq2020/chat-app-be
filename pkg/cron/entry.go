package cron

import "time"

type Entry struct {
	Schedule Schedule
	Job
	Next time.Time
	Prev time.Time
}
