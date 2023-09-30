package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"server/crowdbt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Judge struct {
	Id           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Token        string              `bson:"token" json:"token"`
	Code         string              `bson:"code" json:"code"`
	Name         string              `bson:"name" json:"name"`
	Email        string              `bson:"email" json:"email"`
	Active       bool                `bson:"active" json:"active"`
	LastActivity primitive.DateTime  `bson:"last_activity" json:"last_activity"`
	ReadWelcome  bool                `bson:"read_welcome" json:"read_welcome"`
	Notes        string              `bson:"notes" json:"notes"`
	Votes        int64               `bson:"votes" json:"votes"`
	Next         *primitive.ObjectID `bson:"next" json:"next"`
	Prev         *primitive.ObjectID `bson:"prev" json:"prev"`
	Alpha        float64             `bson:"alpha" json:"alpha"`
	Beta         float64             `bson:"beta" json:"beta"`
	SeenProjects []JudgedProject     `bson:"seen_projects" json:"seen_projects"`
}

type JudgedProject struct {
	ProjectId   primitive.ObjectID `bson:"project_id" json:"project_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Stars       int64              `bson:"stars" json:"stars"`
}

func NewJudge(name string, email string, notes string) *Judge {
	return &Judge{
		Token:        "",
		Code:         fmt.Sprintf("%d", rand.Intn(900000)+100000), // Generates a num between 100000 and 999999
		Name:         name,
		Email:        email,
		Active:       true,
		ReadWelcome:  false,
		Notes:        notes,
		Votes:        0,
		Next:         nil,
		Prev:         nil,
		Alpha:        crowdbt.ALPHA_PRIOR,
		Beta:         crowdbt.BETA_PRIOR,
		LastActivity: primitive.DateTime(0),
		SeenProjects: []JudgedProject{},
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
