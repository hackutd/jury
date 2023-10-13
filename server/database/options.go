package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOptions gets the options from the database
func GetOptions(db *mongo.Database) (*models.Options, error) {
	var options models.Options
	err := db.Collection("options").FindOne(context.Background(), gin.H{}).Decode(&options)

	// If options does not exist, create it
	if err == mongo.ErrNoDocuments {
		options = *models.NewOptions()
		_, err = db.Collection("options").InsertOne(context.Background(), options)
		return &options, err
	}

	return &options, err
}

func UpdateNextTableNum(db *mongo.Database, ctx context.Context, nextTableNum int64) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"next_table_num": nextTableNum}})
	return err
}

func UpdateClock(db *mongo.Database, clock *models.ClockState) error {
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": gin.H{"clock": clock}})
	return err
}
