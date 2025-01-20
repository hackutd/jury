package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Time    int64              `bson:"time" json:"time"`   // First time entry was added to this log
	Count   int64              `bson:"count" json:"count"` // Number of times this log has been updated
	Entries []string           `bson:"entries" json:"entries"`
}

func NewLog() *Log {
	return &Log{
		Time:    GetCurrTime(),
		Count:   0,
		Entries: []string{},
	}
}
