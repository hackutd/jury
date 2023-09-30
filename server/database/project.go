package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertProjects inserts a list of projects into the database
func InsertProjects(db *mongo.Database, projects []*models.Project) error {
	var docs []interface{}
	for _, project := range projects {
		docs = append(docs, project)
	}
	_, err := db.Collection("projects").InsertMany(context.Background(), docs)
	return err
}

// InsertProject inserts a project into the database
func InsertProject(db *mongo.Database, project *models.Project) error {
	_, err := db.Collection("projects").InsertOne(context.Background(), project)
	return err
}

// FindAllProjects returns a list of all projects in the database
func FindAllProjects(db *mongo.Database) ([]*models.Project, error) {
	var projects []*models.Project
	cursor, err := db.Collection("projects").Find(context.Background(), gin.H{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// DeleteProjectById deletes a project from the database by id
func DeleteProjectById(db *mongo.Database, id primitive.ObjectID) error {
	_, err := db.Collection("projects").DeleteOne(context.Background(), gin.H{"_id": id})
	return err
}

func AggregateProjectStats(db *mongo.Database) (*models.ProjectStats, error) {
	// Get the totoal number of projects
	totalProjects, err := db.Collection("projects").EstimatedDocumentCount(context.Background())
	if err != nil {
		return nil, err
	}

	// Get the average votes and average seen using an aggregation pipeline
	cursor, err := db.Collection("projects").Aggregate(context.Background(), []gin.H{
		{"$match": gin.H{"active": true}},
		{"$group": gin.H{
			"_id": nil,
			"avgVotes": gin.H{
				"$avg": "$votes",
			},
			"avgSeen": gin.H{
				"$avg": "$seen",
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var stats models.ProjectStats
	cursor.Next(context.Background())
	err = cursor.Decode(&stats)
	if err != nil {
		return nil, err
	}

	// Set the total number of projects
	stats.Num = totalProjects

	return &stats, nil
}
