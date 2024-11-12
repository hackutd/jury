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

type MoveRequest struct {
	Group int64 `json:"group"`
}

type MoveSelectedRequest struct {
	Items []primitive.ObjectID `json:"items"`
	Group int64                `json:"group"`
}
