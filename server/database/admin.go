package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// AggregateStats aggregates all stats from the database.
func AggregateStats(db *mongo.Database) (*models.Stats, error) {
	// Get the total number of projects and judges
	totalProjects, err := db.Collection("projects").EstimatedDocumentCount(context.Background())
	if err != nil {
		return nil, err
	}
	totalJudges, err := db.Collection("judges").EstimatedDocumentCount(context.Background())
	if err != nil {
		return nil, err
	}

	// Get the average votes, average seen, average mu, and average sigma using an aggregation pipeline
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
			"avgMu": gin.H{
				"$avg": "$mu",
			},
			"avgSigma": gin.H{
				"$avg": "$sigma_sq",
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var stats models.Stats
	cursor.Next(context.Background())
	cursor.Decode(&stats)

	// Set the total number of projects and judges
	stats.Projects = totalProjects
	stats.Judges = totalJudges

	return &stats, nil
}

// DropAll drops the entire database
func DropAll(db *mongo.Database) error {
	// Drop all collections
	var collections = []string{"projects", "judges", "skips", "options", "votes"}
	for _, c := range collections {
		if err := db.Collection(c).Drop(context.Background()); err != nil {
			return err
		}
	}

	return nil
}
