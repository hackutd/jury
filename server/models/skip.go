package models

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// List of valid reasons for skipping a project
var validReasons = []string{"busy", "absent", "cannot-demo", "too-complex", "offensive"}

// Defines an instance where the judge skips a project.
// This can be one of these reasons:
//
//  1. busy: Busy (Being Judged)
//  2. absent: Not Present
//  3. cannot-demo: Cannot Demo Project
//  4. too-complex: Too Complex
//  5. offensive: Offensive Project
//
// With the exception of the 1st reason, all other reasons are grounds for
// flagging, which is defined by the `flag` field.
type Skip struct {
	Id              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ProjectId       *primitive.ObjectID `json:"project_id" bson:"project_id"`
	JudgeId         *primitive.ObjectID `json:"judge_id" bson:"judge_id"`
	Time            primitive.DateTime  `json:"time" bson:"time"`
	ProjectName     string              `json:"project_name" bson:"project_name"`
	ProjectLocation int64               `json:"project_location" bson:"project_location"`
	JudgeName       string              `json:"judge_name" bson:"judge_name"`
	Reason          string              `json:"reason" bson:"reason"`
	Flag            bool                `json:"flag" bson:"flag"`
}

func NewSkip(project *Project, judge *Judge, reason string) (*Skip, error) {
	// Check if the reason is valid
	valid := false
	for _, r := range validReasons {
		if r == reason {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("reason field is invalid: %s", reason)
	}

	// Create the skip object
	return &Skip{
		ProjectId:       &project.Id,
		JudgeId:         &judge.Id,
		Time:            primitive.NewDateTimeFromTime(time.Now()),
		ProjectName:     project.Name,
		ProjectLocation: project.Location,
		JudgeName:       judge.Name,
		Reason:          reason,
		Flag:            reason != "busy",
	}, nil
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
