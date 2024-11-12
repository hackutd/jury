package util

import "go.mongodb.org/mongo-driver/bson/primitive"

// Request structs for the API

type HideRequest struct {
	Hide bool `json:"hide"`
}

type PrioritizeRequest struct {
	Prioritize bool `json:"prioritize"`
}

type HideSelectedRequest struct {
	Items []primitive.ObjectID `json:"items"`
	Hide  bool                 `json:"hide"`
}

type PrioritizeSelectedRequest struct {
	Items      []primitive.ObjectID `json:"items"`
	Prioritize bool                 `json:"prioritize"`
}

type MoveJudgeRequest struct {
	Group int64 `json:"group"`
}

type MoveSelectedJudgesRequest struct {
	Judges []primitive.ObjectID `json:"judges"`
	Group  int64                `json:"group"`
}
