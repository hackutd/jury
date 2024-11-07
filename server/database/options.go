package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOptions gets the options from the database
func GetOptions(db *mongo.Database, ctx context.Context) (*models.Options, error) {
	var options models.Options
	err := db.Collection("options").FindOne(ctx, gin.H{}).Decode(&options)

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
func UpdateCurrTableNum(db *mongo.Database, ctx context.Context, options *models.Options) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"curr_table_num": options.CurrTableNum, "group_table_nums": options.GroupTableNums}})
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
	options, err := GetOptions(db, ctx)
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

func UpdateJudgeTracks(db *mongo.Database, ctx context.Context, judgeTracks bool) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"judge_tracks": judgeTracks}})
	return err
}

func UpdateTracks(db *mongo.Database, ctx context.Context, tracks []string) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"tracks": tracks}})
	return err
}

func UpdateMultiGroup(db *mongo.Database, ctx context.Context, multiGroup bool) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"multi_group": multiGroup}})
	return err
}

func UpdateNumGroups(db *mongo.Database, ctx context.Context, numGroups int64) error {
	// Get options
	options, err := GetOptions(db, ctx)
	if err != nil {
		return err
	}

	// Resize group sizes if necessary
	if numGroups < options.NumGroups {
		options.GroupSizes = options.GroupSizes[:numGroups-1]
	} else if numGroups > options.NumGroups {
		for i := options.NumGroups; i < numGroups-1; i++ {
			options.GroupSizes = append(options.GroupSizes, 30)
		}
	}

	_, err = db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"num_groups": numGroups, "group_sizes": options.GroupSizes}})
	return err
}

// UpdateGroupSizes will update the group sizes based on the given sizes
func UpdateGroupSizes(db *mongo.Database, ctx context.Context, groupSizes []int64) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"group_sizes": groupSizes}})
	return err
}

// UpdateGroupOptions will update the group options based on the given options
func UpdateGroupOptions(db *mongo.Database, ctx context.Context, groupOptions models.OptionalGroupOptions) error {
	update := gin.H{}

	if groupOptions.SwitchingMode != nil {
		update["main_group.switching_mode"] = *groupOptions.SwitchingMode
	}
	if groupOptions.AutoSwitchMethod != nil {
		update["main_group.auto_switch_method"] = *groupOptions.AutoSwitchMethod
	}
	if groupOptions.AutoSwitchCount != nil {
		update["main_group.auto_switch_count"] = *groupOptions.AutoSwitchCount
	}
	if groupOptions.AutoSwitchProp != nil {
		update["main_group.auto_switch_prop"] = *groupOptions.AutoSwitchProp
	}

	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": update})
	return err
}

// IncrementManualSwitches increments the manual switches in the database
func IncrementManualSwitches(db *mongo.Database, ctx context.Context) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$inc": gin.H{"main_group.manual_switches": 1}})
	return err
}
