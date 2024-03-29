package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Options struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ref          int64              `bson:"ref" json:"ref"`
	CurrTableNum int64              `bson:"curr_table_num" json:"curr_table_num"`
	Clock        ClockState         `bson:"clock" json:"clock"`
	Groups       []Group            `bson:"groups" json:"groups"`
	JudgingTimer int64              `bson:"judging_timer" json:"judging_timer"`
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
		JudgingTimer: 300,
		Clock:        *NewClockState(),
	}
}

func NewGroup(start int, end int) *Group {
	return &Group{
		Start: int64(start),
		End:   int64(end),
	}
}

// Type to sort groups by start time
type ByStartTime []Group

func (a ByStartTime) Len() int           { return len(a) }
func (a ByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartTime) Less(i, j int) bool { return a[i].Start < a[j].Start }
