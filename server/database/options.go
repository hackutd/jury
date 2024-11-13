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
func UpdateOptions(db *mongo.Database, ctx context.Context, options *models.OptionalOptions) error {
	update := gin.H{}

	if options.JudgingTimer != nil {
		update["judging_timer"] = *options.JudgingTimer
	}
	if options.MinViews != nil {
		update["min_views"] = *options.MinViews
	}
	if options.ClockSync != nil {
		update["clock_sync"] = *options.ClockSync
	}
	if options.JudgeTracks != nil {
		update["judge_tracks"] = *options.JudgeTracks
	}
	if options.Tracks != nil {
		update["tracks"] = *options.Tracks
	}
	if options.MultiGroup != nil {
		update["multi_group"] = *options.MultiGroup
	}
	if options.NumGroups != nil {
		update["num_groups"] = *options.NumGroups
	}
	if options.GroupSizes != nil {
		update["group_sizes"] = *options.GroupSizes
	}
	if options.SwitchingMode != nil {
		update["switching_mode"] = *options.SwitchingMode
	}
	if options.AutoSwitchProp != nil {
		update["auto_switch_prop"] = *options.AutoSwitchProp
	}
	if options.GroupNames != nil {
		update["group_names"] = *options.GroupNames
	}
	if options.IgnoreTracks != nil {
		update["ignore_tracks"] = *options.IgnoreTracks
	}

	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": update})
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

// UpdateNumGroups will update the number of groups and resize the group sizes if necessary
func UpdateNumGroups(db *mongo.Database, ctx context.Context, numGroups int64) error {
	// Get options
	options, err := GetOptions(db, ctx)
	if err != nil {
		return err
	}

	// Resize group sizes if necessary
	if numGroups < options.NumGroups {
		options.GroupSizes = options.GroupSizes[:numGroups]
	} else if numGroups > options.NumGroups {
		for i := options.NumGroups; i < numGroups; i++ {
			options.GroupSizes = append(options.GroupSizes, 30)
		}
	}

	// Reassign group numbers to all projects
	ReassignAllGroupNums(db, ctx, options)

	// Update options
	return UpdateOptions(db, ctx, &models.OptionalOptions{NumGroups: &numGroups, GroupSizes: &options.GroupSizes})
}

func UpdateGroupSizes(db *mongo.Database, ctx context.Context, groupSizes []int64) error {
	// Get options
	options, err := GetOptions(db, ctx)
	if err != nil {
		return err
	}

	// If group sizes did not change, do nothing
	if len(groupSizes) == len(options.GroupSizes) {
		same := true
		for i := range groupSizes {
			if groupSizes[i] != options.GroupSizes[i] {
				same = false
				break
			}
		}
		if same {
			return nil
		}
	}

	// Reassign group numbers to all projects
	options.GroupSizes = groupSizes
	ReassignAllGroupNums(db, ctx, options)

	// Update options
	return UpdateOptions(db, ctx, &models.OptionalOptions{GroupSizes: &groupSizes})
}

// IncrementManualSwitches increments the manual switches in the database
func IncrementManualSwitches(db *mongo.Database, ctx context.Context) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$inc": gin.H{"manual_switches": 1}})
	return err
}

// UpdateQRCode updates the QR code in the database
func UpdateQRCode(db *mongo.Database, ctx context.Context, qrCode string) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"qr_code": qrCode}})
	return err
}

// UpdateTrackQRCode updates the QR code for a track in the database
func UpdateTrackQRCode(db *mongo.Database, ctx context.Context, track string, qrCode string) error {
	key := "track_qr_codes." + track
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{key: qrCode}})
	return err
}

// UpdateDeliberation updates the deliberation in the database
func UpdateDeliberation(db *mongo.Database, ctx context.Context, deliberation bool) error {
	_, err := db.Collection("options").UpdateOne(ctx, gin.H{}, gin.H{"$set": gin.H{"deliberation": deliberation}})
	return err
}
