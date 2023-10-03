package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Skip struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id"`
	ProjectId *primitive.ObjectID `json:"projectId" bson:"projectId"`
	JudgeId   *primitive.ObjectID `json:"judgeId" bson:"judgeId"`
	Time      primitive.DateTime  `json:"time" bson:"time"`
	Reason    string              `json:"reason" bson:"reason"`
}

func NewSkip(projectId *primitive.ObjectID, judgeId *primitive.ObjectID, reason string) *Skip {
	return &Skip{
		Id:        primitive.NewObjectID(),
		ProjectId: projectId,
		JudgeId:   judgeId,
		Time:      primitive.DateTime(time.Now().UnixNano()),
		Reason:    reason,
	}
}

// Create custom marshal function to change the format of the primitive.DateTime to a unix timestamp
func (s *Skip) MarshalJSON() ([]byte, error) {
	type Alias Skip
	return json.Marshal(&struct {
		*Alias
		Time int64 `json:"time"`
	}{
		Alias: (*Alias)(s),
		Time:  int64(s.Time),
	})
}

// Create custom unmarshal function to change the format of the primitive.DateTime from a unix timestamp
func (s *Skip) UnmarshalJSON(data []byte) error {
	type Alias Skip
	aux := &struct {
		Time int64 `json:"time"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Time = primitive.DateTime(aux.Time)
	return nil
}
