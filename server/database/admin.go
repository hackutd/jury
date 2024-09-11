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
	err = projCursor.Decode(&projAvgSeen)
	if err != nil {
		// This means no documents were found
		if err.Error() == "EOF" {
			projAvgSeen = AvgSeenAgg{AvgSeen: 0}
		} else {
			return nil, err
		}
	}

	// Get the average judge seen using an aggregation pipeline
	judgeCursor, err := db.Collection("judges").Aggregate(context.Background(), []gin.H{
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
	err = judgeCursor.Decode(&judgeAvgSeen)
	if err != nil {
		if err.Error() == "EOF" {
			judgeAvgSeen = AvgSeenAgg{AvgSeen: 0}
		} else {
			return nil, err
		}
	}

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

// UpdateOptions updates the options in the database
func UpdateOptions(db *mongo.Database, options *models.Options) error {
	// Update the options
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": options})
	return err
}

// UpdateCategories updates the categories in the database
func UpdateCategories(db *mongo.Database, categories []string) error {
	// Update the categories
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": gin.H{"categories": categories}})
	return err
}

// UpdateMinViews will update the min views setting
func UpdateMinViews(db *mongo.Database, minViews int) error {
	// Update the min views
	println(minViews)
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": gin.H{"min_views": minViews}})
	return err
}
