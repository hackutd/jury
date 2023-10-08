package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Skip struct {
	Id        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProjectId *primitive.ObjectID `bson:"projectId" json:"projectId"`
	JudgeId   *primitive.ObjectID `bson:"judgeId" json:"judgeId"`
	Time      primitive.DateTime  `bson:"time" json:"time"`
	Reason    string              `bson:"reason" json:"reason"`
}

func NewSkip(projectId *primitive.ObjectID, judgeId *primitive.ObjectID, reason string) *Skip {
	return &Skip{
		ProjectId: projectId,
		JudgeId:   judgeId,
		Time:      primitive.NewDateTimeFromTime(time.Now()),
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
