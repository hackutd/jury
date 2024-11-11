package models

import (
	"encoding/json"

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
	TrackSeen     map[string]int64   `bson:"track_seen" json:"track_seen"`
	Score         int64              `bson:"score" json:"score"`
	Stars         int64              `bson:"stars" json:"stars"`
	TrackStars    map[string]int64   `bson:"track_stars" json:"track_stars"`
	Active        bool               `bson:"active" json:"active"`
	Prioritized   bool               `bson:"prioritized" json:"prioritized"`
	Group         int64              `bson:"group" json:"group"`
	LastActivity  primitive.DateTime `bson:"last_activity" json:"last_activity"`
}

func NewProject(name string, location int64, group int64, description string, url string, tryLink string, videoLink string, challengeList []string) *Project {
	return &Project{
		Name:          name,
		Location:      location,
		Group:         group,
		Description:   description,
		Url:           url,
		TryLink:       tryLink,
		VideoLink:     videoLink,
		ChallengeList: challengeList,
		Seen:          0,
		TrackSeen:     make(map[string]int64),
		Score:         0,
		Stars:         0,
		TrackStars:    make(map[string]int64),
		Active:        true,
		Prioritized:   false,
		LastActivity:  primitive.DateTime(0),
	}
}

func DefaultProject() *Project {
	return NewProject("", 0, 0, "", "", "", "", []string{})
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

// Type to sort by table number
type ByTableNumber []*Project

func (a ByTableNumber) Len() int           { return len(a) }
func (a ByTableNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTableNumber) Less(i, j int) bool { return a[i].Location < a[j].Location }

type HideProjectRequest struct {
	Hide bool `json:"hide"`
}

type PrioritizeProjectRequest struct {
	Prioritize bool `json:"prioritize"`
}
