package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertFlag inserts a skip object into the database
func InsertFlag(db *mongo.Database, ctx context.Context, skip *models.Flag) error {
	_, err := db.Collection("flags").InsertOne(ctx, skip)
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

// FindFlagsByJudge will find the flags that a judge has flagged
func FindFlagsByJudge(db *mongo.Database, judge *models.Judge, ctx mongo.SessionContext) ([]*models.Flag, error) {
	flags := make([]*models.Flag, 0)
	cursor, err := db.Collection("flags").Find(ctx, gin.H{"judge_id": judge.Id})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &flags)
	if err != nil {
		return nil, err
	}
	return flags, nil
}

// GetProjectAbsentCount finds the number of times that a specified project has been skipped
func GetProjectAbsentCount(db *mongo.Database, project *models.Project, ctx mongo.SessionContext) (int, error) {
	flags := make([]*models.Flag, 0)
	cursor, err := db.Collection("flags").Find(ctx, gin.H{"project_id": project.Id, "reason": "absent"})
	if err != nil {
		return 0, err
	}
	err = cursor.All(context.Background(), &flags)
	if err != nil {
		return 0, err
	}
	return len(flags), nil
}

// DeleteAbsentFlags deletes all flags with the reason "absent" for a given project ID
func DeleteAbsentFlags(db *mongo.Database, project *models.Project, ctx mongo.SessionContext) error {
	_, err := db.Collection("flags").DeleteMany(ctx, gin.H{"project_id": project.Id, "reason": "absent"})
	return err
}
