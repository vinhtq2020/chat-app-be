package model

import "time"

type QueryCount struct {
	Query     string    `json:"query" bson:"query"`
	Time      time.Time `json:"time" bson:"time"`
	Prequency int64     `json:"prequency" bson:"prequency"`
}
