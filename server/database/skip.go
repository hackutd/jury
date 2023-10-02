package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertSkip(db *mongo.Database, skip *models.Skip) error {
	_, err := db.Collection("skips").InsertOne(context.Background(), skip)
	return err
}

func FindAllSkips(db *mongo.Database) ([]*models.Skip, error) {
	skips := make([]*models.Skip, 0)
	cursor, err := db.Collection("skips").Find(context.Background(), gin.H{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &skips)
	if err != nil {
		return nil, err
	}
	return skips, nil
}
