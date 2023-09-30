package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertJudge inserts a judge into the database
func InsertJudge(db *mongo.Database, judge *models.Judge) error {
	_, err := db.Collection("judges").InsertOne(context.Background(), judge)
	return err
}

// InsertJudges inserts multiple judges into the database
func InsertJudges(db *mongo.Database, judges []*models.Judge) error {
	var docs []interface{}
	for _, judge := range judges {
		docs = append(docs, judge)
	}
	_, err := db.Collection("judges").InsertMany(context.Background(), docs)
	return err
}

// FindJudgeByToken finds a judge by their token.
// Returns judge as nil if no judge was found.
func FindJudgeByToken(db *mongo.Database, token string) (*models.Judge, error) {
	var judge models.Judge
	err := db.Collection("judges").FindOne(context.Background(), gin.H{"token": token}).Decode(&judge)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &judge, err
}

// FindJudgeByCode finds a judge by their code.
// Returns judge as nil if no judge was found.
func FindJudgeByCode(db *mongo.Database, code string) (*models.Judge, error) {
	var judge models.Judge
	err := db.Collection("judges").FindOne(context.Background(), gin.H{"code": code}).Decode(&judge)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &judge, err
}

// UpdateJudge updates a judge in the database
func UpdateJudge(db *mongo.Database, judge *models.Judge) error {
	_, err := db.Collection("judges").UpdateOne(context.Background(), gin.H{"_id": judge.Id}, gin.H{"$set": judge})
	return err
}

// FindAllJudges returns a list of all judges in the database
func FindAllJudges(db *mongo.Database) ([]*models.Judge, error) {
	var judges []*models.Judge
	cursor, err := db.Collection("judges").Find(context.Background(), gin.H{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var judge models.Judge
		err := cursor.Decode(&judge)
		if err != nil {
			return nil, err
		}
		judges = append(judges, &judge)
	}
	return judges, nil
}

// AggregateJudgeStats aggregates statistics about judges
func AggregateJudgeStats(db *mongo.Database) (*models.JudgeStats, error) {
	// Get the total number of judges
	totalJudges, err := db.Collection("judges").EstimatedDocumentCount(context.Background())
	if err != nil {
		return nil, err
	}

	// Get the total number of active judges and the average number of votes using an aggregation pipeline
	cursor, err := db.Collection("judges").Aggregate(context.Background(), []gin.H{
		{"$match": gin.H{"active": true}},
		{"$group": gin.H{
			"_id": nil,
			"avgVotes": gin.H{
				"$avg": "$votes",
			},
			"numActive": gin.H{
				"$sum": 1,
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var stats models.JudgeStats
	cursor.Next(context.Background())
	err = cursor.Decode(&stats)
	if err != nil {
		return nil, err
	}

	// Set the total number of judges
	stats.Num = totalJudges

	return &stats, nil
}

// DeleteJudgeById deletes a judge from the database by their id
func DeleteJudgeById(db *mongo.Database, id primitive.ObjectID) error {
	_, err := db.Collection("judges").DeleteOne(context.Background(), gin.H{"_id": id})
	return err
}
