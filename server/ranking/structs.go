package ranking

import "go.mongodb.org/mongo-driver/bson/primitive"

type JudgeRanking struct {
	Rankings []primitive.ObjectID `json:"rankings"`
	Unranked []primitive.ObjectID `json:"unranked"`
}
