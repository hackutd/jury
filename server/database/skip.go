package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertSkip inserts a skip object into the database
func InsertSkip(db *mongo.Database, skip *models.Flag) error {
	_, err := db.Collection("skips").InsertOne(context.Background(), skip)
	return err
}

// FindAllSkips returns all skip objects, or all flag objects
func FindAllSkips(db *mongo.Database, flag bool) ([]*models.Flag, error) {
	cond := gin.H{}
	if flag {
		cond["flag"] = true
	}
	cursor, err := db.Collection("skips").Find(context.Background(), cond)
	if err != nil {
		return nil, err
	}

	skips := make([]*models.Flag, 0)
	err = cursor.All(context.Background(), &skips)
	if err != nil {
		return nil, err
	}
	return skips, nil
}
