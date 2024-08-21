package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetComparisons gets the comparisons object from the database
func GetComparisons(db *mongo.Database) (*models.Comparisons, error) {
	var comparisons models.Comparisons
	err := db.Collection("comparisons").FindOne(context.Background(), gin.H{}).Decode(&comparisons)

	// If comparisons object doesn't exit, create it
	if err == mongo.ErrNoDocuments {
		comparisons = *models.NewComparisons()
		_, err = db.Collection("comparisons").InsertOne(context.Background(), comparisons)
		return &comparisons, err
	}

	return &comparisons, err
}

// UpdateComparisons will update the comparisons object in the database
func UpdateComparisons(db *mongo.Database, ctx context.Context, comparisons *models.Comparisons) error {
	_, err := db.Collection("comparisons").UpdateOne(ctx, gin.H{}, gin.H{"$set": comparisons})
	return err
}
