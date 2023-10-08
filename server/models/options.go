package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Options struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ref          int64              `bson:"ref" json:"ref"`
	NextTableNum int64              `bson:"next_table_num" json:"next_table_num"`
	Clock        ClockState         `bson:"clock" json:"clock"`
}

func NewOptions() *Options {
	return &Options{
		Ref:          0,
		NextTableNum: 1,
		Clock:        *NewClockState(),
	}
}
