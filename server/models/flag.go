package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flag struct {
	Id              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ProjectId       *primitive.ObjectID `json:"project_id" bson:"project_id"`
	JudgeId         *primitive.ObjectID `json:"judge_id" bson:"judge_id"`
	Time            primitive.DateTime  `json:"time" bson:"time"`
	ProjectName     string              `json:"project_name" bson:"project_name"`
	ProjectLocation int64               `json:"project_location" bson:"project_location"`
	JudgeName       string              `json:"judge_name" bson:"judge_name"`
	Reason          string              `json:"reason" bson:"reason"`
}

func NewFlag(project *Project, judge *Judge, reason string) *Flag {
	return &Flag{
		ProjectId:       &project.Id,
		JudgeId:         &judge.Id,
		Time:            primitive.NewDateTimeFromTime(time.Now()),
		ProjectName:     project.Name,
		ProjectLocation: project.Location,
		JudgeName:       judge.Name,
		Reason:          reason,
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

// TODO: This is garbage code; just make skip/flag one object smh
func SkipToFlag(skip *Skip) *Flag {
	return &Flag{
		ProjectId:       skip.ProjectId,
		JudgeId:         skip.JudgeId,
		Time:            skip.Time,
		ProjectName:     skip.ProjectName,
		ProjectLocation: skip.ProjectLocation,
		JudgeName:       skip.JudgeName,
		Reason:          skip.Reason,
	}
}
