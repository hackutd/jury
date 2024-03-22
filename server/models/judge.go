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
	ReadWelcome  bool                 `bson:"read_welcome" json:"read_welcome"`
	Notes        string               `bson:"notes" json:"notes"`
	Current      *primitive.ObjectID  `bson:"current" json:"current"`
	Seen         int64                `bson:"seen" json:"seen"`
	SeenProjects []JudgedProject      `bson:"seen_projects" json:"seen_projects"`
	Rankings     []primitive.ObjectID `bson:"rankings" json:"rankings"`
	LastActivity primitive.DateTime   `bson:"last_activity" json:"last_activity"`
}

type JudgedProject struct {
	ProjectId   primitive.ObjectID `bson:"project_id" json:"project_id"`
	Categories  map[string]int     `bson:"categories" json:"categories"`
	Notes       string             `bson:"notes" json:"notes"`
	Name        string             `bson:"name" json:"name"`
	Location    int64              `bson:"location" json:"location"`
	Description string             `bson:"description" json:"description"`
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
		Current:      nil,
		Seen:         0,
		SeenProjects: []JudgedProject{},
		Rankings:     []primitive.ObjectID{},
		LastActivity: primitive.DateTime(0),
	}
}

func JudgeProjectFromProject(project *Project, categories map[string]int) *JudgedProject {
	return &JudgedProject{
		ProjectId:   project.Id,
		Categories:  categories,
		Name:        project.Name,
		Location:    project.Location,
		Description: project.Description,
		Notes:       "",
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
