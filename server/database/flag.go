package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertFlag(db *mongo.Database, flag *models.Flag) error {
	_, err := db.Collection("flags").InsertOne(context.Background(), flag)
	return err
}

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
