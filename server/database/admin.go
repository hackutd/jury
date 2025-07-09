package database

import (
	"context"
	"fmt"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// DropProjects removes all projects
func DropProjects(db *mongo.Database) error {
	fmt.Println("??????")
	err := db.Collection("projects").Drop(context.Background())
	return err
}

// DropJudges removes all judges
func DropJudges(db *mongo.Database) error {
	err := db.Collection("judges").Drop(context.Background())
	return err
}

// DropJudgingData removes all fields from judges/projects pertaining to judging
func DropJudgingData(db *mongo.Database) error {
	_, err := db.Collection("judges").UpdateMany(
		context.Background(),
		gin.H{},
		gin.H{"$set": gin.H{
			"active":        true,
			"current":       nil,
			"last_location": -1,
			"seen":          0,
			"group_seen":    0,
			"read_welcome":  false,
			"seen_projects": []models.JudgedProject{},
			"rankings":      []primitive.ObjectID{},
			"rankings_agg":  []models.AggRanking{},
		}},
	)
	if err != nil {
		return err
	}

	_, err = db.Collection("projects").UpdateMany(
		context.Background(),
		gin.H{},
		gin.H{"$set": gin.H{
			"active":      true,
			"seen":        0,
			"track_seen":  map[string]int64{},
			"prioritized": false,
		}},
	)
	if err != nil {
		return err
	}

	err = db.Collection("flags").Drop(context.Background())
	if err != nil {
		return err
	}

	_, err = db.Collection("options").UpdateOne(
		context.Background(),
		gin.H{"ref": 0},
		gin.H{"$set": gin.H{
			"clock":           *models.NewClockState(),
			"deliberation":    false,
			"manual_switches": 0,
			"qr_code":         "",
			"track_qr_codes":  make(map[string]string),
		}},
	)
	return err
}

// DropRankings removes the rankings from every judge
func DropRankings(db *mongo.Database) error {
	_, err := db.Collection("judges").UpdateMany(
		context.Background(),
		gin.H{},
		gin.H{"$set": gin.H{
			"rankings":                  []primitive.ObjectID{},
			"rankings_agg":              []models.AggRanking{},
			"seen_projects.$[].starred": false,
		}},
	)
	return err
}

// UpdateMinViews will update the min views setting
func UpdateMinViews(db *mongo.Database, minViews int) error {
	// Update the min views
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": gin.H{"min_views": minViews}})
	return err
}
