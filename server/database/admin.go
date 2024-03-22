package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AvgSeenAgg struct {
	AvgSeen float64 `bson:"avgSeen"`
}

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

	// Get the average project seen using an aggregation pipeline
	projCursor, err := db.Collection("projects").Aggregate(context.Background(), []gin.H{
		{"$match": gin.H{"active": true}},
		{"$group": gin.H{
			"_id": nil,
			"avgSeen": gin.H{
				"$avg": "$seen",
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var projAvgSeen AvgSeenAgg
	projCursor.Next(context.Background())
	projCursor.Decode(&projAvgSeen)

	// Get the average judge seen using an aggregation pipeline
	judgeCursor, err := db.Collection("judge").Aggregate(context.Background(), []gin.H{
		{"$match": gin.H{"active": true}},
		{"$group": gin.H{
			"_id": nil,
			"avgSeen": gin.H{
				"$avg": "$seen",
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var judgeAvgSeen AvgSeenAgg
	judgeCursor.Next(context.Background())
	judgeCursor.Decode(&judgeAvgSeen)

	// Create the stats object
	var stats models.Stats

	// Set stats from aggregations
	stats.Projects = totalProjects
	stats.Judges = totalJudges
	stats.AvgProjectSeen = projAvgSeen.AvgSeen
	stats.AvgJudgeSeen = judgeAvgSeen.AvgSeen

	return &stats, nil
}

// DropAll drops the entire database
func DropAll(db *mongo.Database) error {
	// Drop all collections
	var collections = []string{"projects", "judges", "flags", "options"}
	for _, c := range collections {
		if err := db.Collection(c).Drop(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

func UpdateOptions(db *mongo.Database, options *models.Options) error {
	// Update the options
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": options})
	return err
}
