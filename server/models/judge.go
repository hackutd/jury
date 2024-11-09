package models

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Judge struct {
	Id           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Token        string               `bson:"token" json:"token"`
	Code         string               `bson:"code" json:"code"`
	Name         string               `bson:"name" json:"name"`
	Email        string               `bson:"email" json:"email"`
	Active       bool                 `bson:"active" json:"active"`
	Track        string               `bson:"track" json:"track"`
	Group        int64                `bson:"group" json:"group"`
	ReadWelcome  bool                 `bson:"read_welcome" json:"read_welcome"`
	Notes        string               `bson:"notes" json:"notes"`
	Current      *primitive.ObjectID  `bson:"current" json:"current"`
	LastLocation int64                `bson:"last_location" json:"last_location"`
	Seen         int64                `bson:"seen" json:"seen"`
	GroupSeen    int64                `bson:"group_seen" json:"group_seen"` // Projects seen in the group
	SeenProjects []JudgedProject      `bson:"seen_projects" json:"seen_projects"`
	Rankings     []primitive.ObjectID `bson:"rankings" json:"rankings"`
	LastActivity primitive.DateTime   `bson:"last_activity" json:"last_activity"`
}

type JudgedProject struct {
	ProjectId   primitive.ObjectID `bson:"project_id" json:"project_id"`
	Starred     bool               `bson:"starred" json:"starred"`
	Notes       string             `bson:"notes" json:"notes"`
	Name        string             `bson:"name" json:"name"`
	Location    int64              `bson:"location" json:"location"`
	Description string             `bson:"description" json:"description"`
}

func NewJudge(name string, email string, track string, notes string, group int64) *Judge {
	return &Judge{
		Token:        "",
		Code:         fmt.Sprintf("%d", rand.Intn(900000)+100000), // Generates a num between 100000 and 999999
		Name:         name,
		Email:        email,
		Active:       true,
		Group:        group,
		Track:        track,
		ReadWelcome:  false,
		Notes:        notes,
		Current:      nil,
		LastLocation: -1,
		Seen:         0,
		GroupSeen:    0,
		SeenProjects: []JudgedProject{},
		Rankings:     []primitive.ObjectID{},
		LastActivity: primitive.DateTime(0),
	}
}

func JudgeProjectFromProject(project *Project, notes string, starred bool) *JudgedProject {
	return &JudgedProject{
		ProjectId:   project.Id,
		Name:        project.Name,
		Location:    project.Location,
		Description: project.Description,
		Notes:       notes,
		Starred:     starred,
	}
}

// Create custom marshal function to change the format of the primitive.DateTime to a unix timestamp
func (j *Judge) MarshalJSON() ([]byte, error) {
	type Alias Judge
	return json.Marshal(&struct {
		*Alias
		LastActivity int64 `json:"last_activity"`
	}{
		Alias:        (*Alias)(j),
		LastActivity: int64(j.LastActivity),
	})
}

// Create custom unmarshal function to change the format of the primitive.DateTime from a unix timestamp
func (j *Judge) UnmarshalJSON(data []byte) error {
	type Alias Judge
	aux := &struct {
		LastActivity int64 `json:"last_activity"`
		*Alias
	}{
		Alias: (*Alias)(j),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	j.LastActivity = primitive.DateTime(aux.LastActivity)
	return nil
}
