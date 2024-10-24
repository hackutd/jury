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

// UpdateOptions updates the options in the database
func UpdateOptions(db *mongo.Database, options *models.Options) error {
	// Update the options
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": options})
	return err
}

// UpdateCurrTableNum updates the current table number in the database
func UpdateCurrTableNum(db *mongo.Database, ctx context.Context, currTableNum int64) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"curr_table_num": currTableNum}})
	return err
}

// UpdateClockSync updates the clock sync in the database
func UpdateClockSync(db *mongo.Database, ctx context.Context, clockSync bool) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"clock_sync": clockSync}})
	return err
}

// UpdateClockConditional updates the clock in the database if clock sync is enabled
func UpdateClockConditional(db *mongo.Database, ctx context.Context, clock *models.ClockState) error {
	// Get options
	options, err := GetOptions(db)
	if err != nil {
		return err
	}

	// If clock sync is not enabled, don't sync the clock
	if !options.ClockSync {
		return nil
	}

	_, err = db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"clock": clock}})
	return err
}

// UpdateClock updates the clock in the database
func UpdateClock(db *mongo.Database, clock *models.ClockState) error {
	_, err := db.Collection("options").UpdateOne(context.Background(), gin.H{}, gin.H{"$set": gin.H{"clock": clock}})
	return err
}

// GetCategories gets the categories from the database
func GetCategories(db *mongo.Database) ([]string, error) {
	var options models.Options
	err := db.Collection("options").FindOne(context.Background(), gin.H{}).Decode(&options)
	return options.Categories, err
}

// GetMinViews gets the minimum views option from the database
func GetMinViews(db *mongo.Database) (int64, error) {
	var options models.Options
	err := db.Collection("options").FindOne(context.Background(), gin.H{}).Decode(&options)
	return options.MinViews, err
}
