package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Options struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ref          int64              `bson:"ref" json:"ref"`
	CurrTableNum int64              `bson:"curr_table_num" json:"curr_table_num"`
	Clock        ClockState         `bson:"clock" json:"clock"`
	Groups       []Group            `bson:"groups" json:"groups"`
	UseGroups    bool               `bson:"use_groups" json:"use_groups"`
}

type Group struct {
	Start int64 `bson:"start" json:"start"`
	End   int64 `bson:"end" json:"end"`
}

func NewOptions() *Options {
	return &Options{
		Ref:          0,
		CurrTableNum: 0,
		Clock:        *NewClockState(),
	}
}

func NewGroup(start int, end int) *Group {
	return &Group{
		Start: int64(start),
		End:   int64(end),
	}
}
