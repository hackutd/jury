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

type MoveGroupRequest struct {
	Group int64 `json:"group"`
}

type MoveRequest struct {
	Location int64 `json:"location"`
}

type MoveSelectedRequest struct {
	Items []primitive.ObjectID `json:"items"`
	Group int64                `json:"group"`
}
