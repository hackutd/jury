package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comparisons struct {
	Id  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ref int64              `bson:"ref" json:"ref"`
	Arr [][]int32          `bson:"arr" json:"arr"`
}

func NewComparisons() *Comparisons {
	return &Comparisons{
		Ref: 0,
		Arr: make([][]int32, 0),
	}
}
