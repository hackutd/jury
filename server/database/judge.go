package database

import (
	"context"
	"errors"
	"server/models"
	"server/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertJudge inserts a judge into the database
func InsertJudge(db *mongo.Database, ctx context.Context, judge *models.Judge) error {
	_, err := db.Collection("judges").InsertOne(ctx, judge)
	return err
}

// InsertJudges inserts multiple judges into the database
func InsertJudges(db *mongo.Database, ctx context.Context, judges []*models.Judge) error {
	var docs []any
	for _, judge := range judges {
		docs = append(docs, judge)
	}
	_, err := db.Collection("judges").InsertMany(ctx, docs)
	return err
}

// FindJudgeByToken finds a judge by their token.
// Returns judge as nil if no judge was found.
func FindJudgeByToken(db *mongo.Database, token string) (*models.Judge, error) {
	var judge models.Judge
	err := db.Collection("judges").FindOne(context.Background(), gin.H{"token": token}).Decode(&judge)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &judge, err
}

// FindJudgeByCode finds a judge by their code.
// Returns judge as nil if no judge was found.
func FindJudgeByCode(db *mongo.Database, ctx context.Context, code string) (*models.Judge, error) {
	var judge models.Judge
	err := db.Collection("judges").FindOne(ctx, gin.H{"code": code}).Decode(&judge)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &judge, err
}

func FindJudgesByTrack(db *mongo.Database, ctx context.Context, track string) ([]*models.Judge, error) {
	judges := make([]*models.Judge, 0)
	cursor, err := db.Collection("judges").Find(ctx, gin.H{"track": track})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &judges)
	if err != nil {
		return nil, err
	}
	return judges, nil
}

// UpdateJudgeToken updates the token of a judge
func UpdateJudgeToken(db *mongo.Database, ctx context.Context, judgeId *primitive.ObjectID, token string) error {
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judgeId}, gin.H{"$set": gin.H{"token": token, "last_activity": util.Now()}})
	return err
}

// UpdateJudgeReadWelcome updates the read_welcome field of a judge to true
func UpdateJudgeReadWelcome(db *mongo.Database, ctx context.Context, judgeId *primitive.ObjectID) error {
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judgeId}, gin.H{"$set": gin.H{"read_welcome": true, "last_activity": util.Now()}})
	return err
}

// UpdateJudgeGroups updates the group and group_seen fields of multiple judges
func UpdateJudgeGroups(db *mongo.Database, sc mongo.SessionContext, judges []*models.Judge) error {
	models := make([]mongo.WriteModel, 0, len(judges))
	for _, judge := range judges {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(gin.H{"_id": judge.Id}).SetUpdate(gin.H{"$set": gin.H{"group": judge.Group, "group_seen": judge.GroupSeen}}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := db.Collection("judges").BulkWrite(sc, models, opts)
	return err
}

// FindAllJudges returns a list of all judges in the database
func FindAllJudges(db *mongo.Database, ctx context.Context) ([]*models.Judge, error) {
	judges := make([]*models.Judge, 0)
	cursor, err := db.Collection("judges").Find(context.Background(), gin.H{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var judge models.Judge
		err := cursor.Decode(&judge)
		if err != nil {
			return nil, err
		}
		judges = append(judges, &judge)
	}
	return judges, nil
}

// AggregateJudgeStats aggregates statistics about judges
func AggregateJudgeStats(db *mongo.Database) (*models.JudgeStats, error) {
	// Get the total number of judges
	totalJudges, err := db.Collection("judges").EstimatedDocumentCount(context.Background())
	if err != nil {
		return nil, err
	}

	// Get the total number of active judges and the average number of votes using an aggregation pipeline
	cursor, err := db.Collection("judges").Aggregate(context.Background(), []gin.H{
		{"$match": gin.H{"active": true}},
		{"$group": gin.H{
			"_id": nil,
			"avgSeen": gin.H{
				"$avg": "$seen",
			},
			"numActive": gin.H{
				"$sum": 1,
			},
		}},
	})
	if err != nil {
		return nil, err
	}

	// Get the first document from the cursor
	var stats models.JudgeStats
	cursor.Next(context.Background())
	err = cursor.Decode(&stats)
	if err != nil {
		if err.Error() == "EOF" {
			stats = models.JudgeStats{Num: 0, AvgSeen: 0, NumActive: 0}
		} else {
			return nil, err
		}
	}

	// Set the total number of judges
	stats.Num = totalJudges

	return &stats, nil
}

// DeleteJudgeById deletes a judge from the database by their id
func DeleteJudgeById(db *mongo.Database, id primitive.ObjectID) error {
	_, err := db.Collection("judges").DeleteOne(context.Background(), gin.H{"_id": id})
	return err
}

// UpdateAfterSeen updates the judge's seen projects and increments the seen count.
// This should be run in a transaction.
func UpdateAfterSeen(db *mongo.Database, ctx context.Context, judge *models.Judge, seenProject *models.JudgedProject) error {
	// Update the judge's seen projects
	_, err := db.Collection("judges").UpdateOne(
		ctx,
		gin.H{"_id": judge.Id},
		gin.H{
			"$push": gin.H{"seen_projects": seenProject},
			"$inc":  gin.H{"seen": 1},
			"$set":  gin.H{"current": nil, "group": judge.Group, "group_seen": judge.GroupSeen, "last_activity": util.Now()},
		},
	)
	if err != nil {
		return errors.New("error updating judge: " + err.Error())
	}

	star := 0
	if seenProject.Starred {
		star = 1
	}

	incData := gin.H{}
	if judge.Track != "" {
		incData["track_seen."+judge.Track] = 1
		incData["track_stars."+judge.Track] = star
	} else {
		incData["seen"] = 1
		incData["stars"] = star
	}

	// Update the project's seen count
	_, err = db.Collection("projects").UpdateOne(
		ctx,
		gin.H{"_id": seenProject.ProjectId},
		gin.H{"$inc": incData, "$set": gin.H{"last_activity": util.Now()}},
	)
	if err != nil {
		return errors.New("error updating project: " + err.Error())
	}

	return nil
}

// SetJudgeActive sets the active field of a judge
func SetJudgeActive(db *mongo.Database, id *primitive.ObjectID, active bool) error {
	_, err := db.Collection("judges").UpdateOne(
		context.Background(),
		gin.H{"_id": id},
		gin.H{"$set": gin.H{"active": active}},
	)
	return err
}

// SetJudgesActive sets the active field of multiple judges
func SetJudgesActive(db *mongo.Database, ids []primitive.ObjectID, active bool) error {
	_, err := db.Collection("judges").UpdateMany(
		context.Background(),
		gin.H{"_id": gin.H{"$in": ids}},
		gin.H{"$set": gin.H{"active": active}},
	)
	return err
}

// UpdateJudgeBasicInfo updates the basic info of a judge (name, email, notes)
func UpdateJudgeBasicInfo(db *mongo.Database, judgeId *primitive.ObjectID, addRequest *models.AddJudgeRequest) error {
	_, err := db.Collection("judges").UpdateOne(
		context.Background(),
		gin.H{"_id": judgeId},
		gin.H{"$set": gin.H{"name": addRequest.Name, "email": addRequest.Email, "notes": addRequest.Notes}},
	)
	return err
}

// UpdateJudgeRanking updates the judge's ranking array
func UpdateJudgeRanking(db *mongo.Database, ctx context.Context, judge *models.Judge, rankings []primitive.ObjectID) error {
	_, err := db.Collection("judges").UpdateOne(
		ctx,
		gin.H{"_id": judge.Id},
		gin.H{"$set": gin.H{"rankings": rankings}},
	)
	return err
}

// UpdateJudgeSeenProjects updates the seen projects of a judge
func UpdateJudgeSeenProjects(db *mongo.Database, ctx context.Context, judge *models.Judge) error {
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judge.Id}, gin.H{"$set": gin.H{"seen_projects": judge.SeenProjects}})
	return err
}

// UpdateJudgeNotes updates the notes of a single seen project for a judge
func UpdateJudgeNotes(db *mongo.Database, ctx context.Context, judgeId *primitive.ObjectID, projIndex int, notes string) error {
	notesKey := "seen_projects." + strconv.Itoa(projIndex) + ".notes"
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judgeId}, gin.H{"$set": gin.H{notesKey: notes}})
	return err
}

// UpdateJudgeStars updates the star value of a single seen project for a judge
func UpdateJudgeStars(db *mongo.Database, ctx context.Context, judgeId primitive.ObjectID, projIndex int, starred bool) error {
	starredKey := "seen_projects." + strconv.Itoa(projIndex) + ".starred"
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judgeId}, gin.H{"$set": gin.H{starredKey: starred}})
	return err

}

// Reset list of projects that judge skipped due to busy
func ResetBusyProjectListForJudge(db *mongo.Database, ctx context.Context, judge *models.Judge) error {
	_, err := db.Collection("flags").DeleteMany(ctx, gin.H{
		"judge_id": judge.Id,
		"reason":   "busy",
	})
	if err != nil {
		return err
	}
	return nil
}

// GetMinJudgeGroup returns the group with the fewest judges
func GetMinJudgeGroup(db *mongo.Database, ctx context.Context) (int64, error) {
	pipe := []gin.H{{"$match": gin.H{
		"active": true,
		"track":  "",
	}}, {"$group": gin.H{
		"_id":   "$group",
		"count": gin.H{"$sum": 1},
	}}}

	cursor, err := db.Collection("judges").Aggregate(ctx, pipe)
	if err != nil {
		return -1, err
	}

	groupCounts := make(map[int64]int64)
	for cursor.Next(ctx) {
		var result map[string]any
		err := cursor.Decode(&result)
		if err != nil {
			return -1, err
		}
		groupCounts[result["_id"].(int64)] = int64(result["count"].(int32))
	}

	// Sort by count in ascending order
	minGroups := util.SortMapByValue(groupCounts, false)

	// Get options
	options, err := GetOptions(db, ctx)
	if err != nil {
		return -1, err
	}

	// Find all groups that aren't in minGroups
	group := int64(-1)
	for i := int64(0); i < options.NumGroups; i++ {
		if _, ok := groupCounts[i]; !ok {
			group = i
			break
		}
	}

	// If there are no groups that aren't in minGroups, use the first group
	if group == -1 {
		group = minGroups[0]
	}

	return group, nil
}

// GetNextNJudgeGroups returns the next n groups to assign new judges to
// The groups are chosen based on the current number of judges in each group,
// with the goal to balance the number of judges in each group
func GetNextNJudgeGroups(db *mongo.Database, ctx context.Context, n int, reset bool) ([]int64, error) { // Get options
	// Get options
	options, err := GetOptions(db, ctx)
	if err != nil {
		return nil, err
	}

	// Get the current number of judges in each group
	var groupCounts map[int64]int64
	if reset {
		// Simply set all group counts to 0
		groupCounts = make(map[int64]int64, options.NumGroups)
	} else {
		groupCounts = make(map[int64]int64)
		cursor, err := db.Collection("judges").Aggregate(ctx, []gin.H{
			{"$match": gin.H{
				"active": true,
				"track":  "",
			}},
			{"$group": gin.H{
				"_id":   "$group",
				"count": gin.H{"$sum": 1},
			}},
		})
		if err != nil {
			return nil, err
		}

		for cursor.Next(ctx) {
			var result map[string]any
			err := cursor.Decode(&result)
			if err != nil {
				return nil, err
			}
			groupCounts[result["_id"].(int64)] = int64(result["count"].(int32))
		}

		// Find all groups that aren't in groups
		for i := int64(0); i < options.NumGroups; i++ {
			if _, ok := groupCounts[i]; !ok {
				groupCounts[i] = 0
			}
		}
	}

	// Keep looping until we have enough groups
	groupList := make([]int64, n)
	for i := 0; i < n; i++ {
		min := groupCounts[0]
		minGroup := int64(0)
		for group, count := range groupCounts {
			if count < min {
				min = count
				minGroup = group
			}
		}
		groupList[i] = minGroup
		groupCounts[minGroup] -= 1
	}

	return groupList, nil
}

// PutJudgesInGroups assigns judges to groups
func PutJudgesInGroups(db *mongo.Database) error {
	return WithTransaction(db, func(sc mongo.SessionContext) error {
		// Get all judges
		judges, err := FindJudgesByTrack(db, sc, "")
		if err != nil {
			return err
		}

		// Get the next n groups to assign judges to
		groups, err := GetNextNJudgeGroups(db, sc, len(judges), true)
		if err != nil {
			return err
		}

		// Assign each judge to a group
		for i, judge := range judges {
			judge.Group = groups[i]
			judge.GroupSeen = 0
		}

		// Update the judges in the database
		err = UpdateJudgeGroups(db, sc, judges)
		return err
	})
}

// SetJudgeGroup sets the group of a judge
func SetJudgeGroup(db *mongo.Database, ctx context.Context, judgeId *primitive.ObjectID, group int64) error {
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judgeId}, gin.H{"$set": gin.H{"group": group, "group_seen": 0}})
	return err
}

// SetJudgesGroup sets the group of multiple judges
func SetJudgesGroup(db *mongo.Database, ctx context.Context, judgeIds []primitive.ObjectID, group int64) error {
	_, err := db.Collection("judges").UpdateMany(ctx, gin.H{"_id": gin.H{"$in": judgeIds}}, gin.H{"$set": gin.H{"group": group, "group_seen": 0}})
	return err
}
