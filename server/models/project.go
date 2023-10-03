package models

import (
	"encoding/json"
	"server/crowdbt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Location      int64              `bson:"location" json:"location"`
	Description   string             `bson:"description" json:"description"`
	Url           string             `bson:"url" json:"url"`
	TryLink       string             `bson:"try_link" json:"try_link"`
	VideoLink     string             `bson:"video_link" json:"video_link"`
	ChallengeList []string           `bson:"challenge_list" json:"challenge_list"`
	Seen          int64              `bson:"seen" json:"seen"`
	Votes         int64              `bson:"votes" json:"votes"`
	Mu            float64            `bson:"mu" json:"mu"`
	SigmaSq       float64            `bson:"sigma_sq" json:"sigma_sq"`
	Active        bool               `bson:"active" json:"active"`
	Prioritized   bool               `bson:"prioritized" json:"prioritized"`
	LastActivity  primitive.DateTime `bson:"last_activity" json:"last_activity"`
}

func NewProject(name string, location int64, description string, url string, tryLink string, videoLink string, challengeList []string) *Project {
	return &Project{
		Name:          name,
		Location:      location,
		Description:   description,
		Url:           url,
		TryLink:       tryLink,
		VideoLink:     videoLink,
		ChallengeList: challengeList,
		Seen:          0,
		Votes:         0,
		Mu:            crowdbt.MU_PRIOR,
		SigmaSq:       crowdbt.SIGMA_SQ_PRIOR,
		Active:        true,
		Prioritized:   false,
		LastActivity:  primitive.DateTime(0),
	}
}

func DefaultProject() *Project {
	return NewProject("", 0, "", "", "", "", []string{})
}

// Create custom marshal function to change the format of the primitive.DateTime to a unix timestamp
func (p *Project) MarshalJSON() ([]byte, error) {
	type Alias Project
	return json.Marshal(&struct {
		*Alias
		LastActivity int64 `json:"last_activity"`
	}{
		Alias:        (*Alias)(p),
		LastActivity: int64(p.LastActivity),
	})
}

// Create custom unmarshal function to change the format of the primitive.DateTime from a unix timestamp
func (p *Project) UnmarshalJSON(data []byte) error {
	type Alias Project
	aux := &struct {
		LastActivity int64 `json:"last_activity"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.LastActivity = primitive.DateTime(aux.LastActivity)
	return nil
}
