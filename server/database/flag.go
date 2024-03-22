package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertFlag inserts a skip object into the database
func InsertFlag(db *mongo.Database, skip *models.Flag) error {
	_, err := db.Collection("flags").InsertOne(context.Background(), skip)
	return err
}

// FindAllFlags returns all flag objects
func FindAllFlags(db *mongo.Database) ([]*models.Flag, error) {
	flags := make([]*models.Flag, 0)
	cursor, err := db.Collection("flags").Find(context.Background(), gin.H{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &flags)
	if err != nil {
		return nil, err
	}
	return flags, nil
}

func FindFlagsByJudge(db *mongo.Database, judge *models.Judge) ([]*models.Flag, error) {
	flags := make([]*models.Flag, 0)
	cursor, err := db.Collection("flags").Find(context.Background(), gin.H{"judge_id": judge.Id})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &flags)
	if err != nil {
		return nil, err
	}
	return flags, nil
}
