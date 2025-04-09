package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsAgg struct {
	AvgSeen float64 `bson:"avgSeen"`
	Count   int64   `bson:"count"`
}

// AggregateStats aggregates all stats from the database.
func AggregateStats(db *mongo.Database, track string) (*models.Stats, error) {
	// Get the average project seen using an aggregation pipeline
	matchObj := gin.H{"active": true}
	if track != "" {
		matchObj["challenge_list"] = track
	}
	projCursor, err := db.Collection("projects").Aggregate(context.Background(), []gin.H{
		{"$match": matchObj},
		{"$group": gin.H{
			"_id": nil,
			"avgSeen": gin.H{
				"$avg": "$seen",
			},
			"count": gin.H{
				"$sum": 1,
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var projStatAgg StatsAgg
	projCursor.Next(context.Background())
	err = projCursor.Decode(&projStatAgg)
	if err != nil {
		// This means no documents were found
		if err.Error() == "EOF" {
			projStatAgg = StatsAgg{AvgSeen: 0, Count: 0}
		} else {
			return nil, err
		}
	}

	// Get the average judge seen using an aggregation pipeline
	matchObj = gin.H{"active": true, "track": ""}
	if track != "" {
		matchObj["track"] = track
	}
	judgeCursor, err := db.Collection("judges").Aggregate(context.Background(), []gin.H{
		{"$match": matchObj},
		{"$group": gin.H{
			"_id": nil,
			"avgSeen": gin.H{
				"$avg": "$seen",
			},
			"count": gin.H{
				"$sum": 1,
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var judgeStatAgg StatsAgg
	judgeCursor.Next(context.Background())
	err = judgeCursor.Decode(&judgeStatAgg)
	if err != nil {
		if err.Error() == "EOF" {
			judgeStatAgg = StatsAgg{AvgSeen: 0, Count: 0}
		} else {
			return nil, err
		}
	}

	// Create the stats object
	var stats models.Stats

	// Set stats from aggregations
	// stats.Projects = totalProjects
	// stats.Judges = totalJudges
	stats.AvgProjectSeen = projStatAgg.AvgSeen
	stats.AvgJudgeSeen = judgeStatAgg.AvgSeen

	// Get count if track
	// if track != "" {
	stats.Projects = projStatAgg.Count
	stats.Judges = judgeStatAgg.Count
	// }

	return &stats, nil
}

// DropAll drops the entire database
func DropAll(db *mongo.Database) error {
	// Drop all collections
	var collections = []string{"projects", "judges", "flags", "options", "logs"}
	for _, c := range collections {
		if err := db.Collection(c).Drop(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMinViews will update the min views setting
func UpdateMinViews(db *mongo.Database, minViews int) error {
	// Update the min views
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": gin.H{"min_views": minViews}})
	return err
}
