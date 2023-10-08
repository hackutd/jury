package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flag struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ProjectId *primitive.ObjectID `json:"projectId" bson:"projectId"`
	JudgeId   *primitive.ObjectID `json:"judgeId" bson:"judgeId"`
	Time      primitive.DateTime  `json:"time" bson:"time"`
	Reason    string              `json:"reason" bson:"reason"`
}

func NewFlag(projectId *primitive.ObjectID, judgeId *primitive.ObjectID, reason string) *Flag {
	return &Flag{
		ProjectId: projectId,
		JudgeId:   judgeId,
		Time:      primitive.NewDateTimeFromTime(time.Now()),
		Reason:    reason,
	}
}

// Create custom marshal function to change the format of the primitive.DateTime to a unix timestamp
func (f *Flag) MarshalJSON() ([]byte, error) {
	type Alias Flag
	return json.Marshal(&struct {
		*Alias
		Time int64 `json:"time"`
	}{
		Alias: (*Alias)(f),
		Time:  int64(f.Time),
	})
}

// Create custom unmarshal function to change the format of the primitive.DateTime from a unix timestamp
func (f *Flag) UnmarshalJSON(data []byte) error {
	type Alias Flag
	aux := &struct {
		Time int64 `json:"time"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.Time = primitive.DateTime(aux.Time)
	return nil
}
