package database

import (
	"context"
	"server/models"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateJudgeLastActivity to the current time
// TODO: Do we actually need this???
func UpdateJudgeLastActivity(db *mongo.Database, ctx context.Context, id *primitive.ObjectID) error {
	lastActivity := util.Now()
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": id}, gin.H{"$set": gin.H{"last_activity": lastActivity}})
	return err
}

// InsertJudge inserts a judge into the database
func InsertJudge(db *mongo.Database, judge *models.Judge) error {
	_, err := db.Collection("judges").InsertOne(context.Background(), judge)
	return err
}

// InsertJudges inserts multiple judges into the database
func InsertJudges(db *mongo.Database, judges []*models.Judge) error {
	var docs []interface{}
	for _, judge := range judges {
		docs = append(docs, judge)
	}
	_, err := db.Collection("judges").InsertMany(context.Background(), docs)
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
func FindJudgeByCode(db *mongo.Database, code string) (*models.Judge, error) {
	var judge models.Judge
	err := db.Collection("judges").FindOne(context.Background(), gin.H{"code": code}).Decode(&judge)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &judge, err
}

// UpdateJudge updates a judge in the database
func UpdateJudge(db *mongo.Database, judge *models.Judge) error {
	judge.LastActivity = util.Now()
	_, err := db.Collection("judges").UpdateOne(context.Background(), gin.H{"_id": judge.Id}, gin.H{"$set": judge})
	return err
}

// FindAllJudges returns a list of all judges in the database
func FindAllJudges(db *mongo.Database) ([]*models.Judge, error) {
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

// UpdateAfterSeen updates the judge's seen projects and increments the seen count
func UpdateAfterSeen(db *mongo.Database, judge *models.Judge, seenProject *models.JudgedProject) error {
	// Update the judge's seen projects
	_, err := db.Collection("judges").UpdateOne(
		context.Background(),
		gin.H{"_id": judge.Id},
		gin.H{
			"$push": gin.H{"seen_projects": seenProject},
			"$inc":  gin.H{"seen": 1},
			"$set":  gin.H{"current": nil, "last_activity": util.Now()},
		},
	)
	return err
}

// SetJudgeHidden sets the active field of a judge
func SetJudgeHidden(db *mongo.Database, id *primitive.ObjectID, hidden bool) error {
	_, err := db.Collection("judges").UpdateOne(
		context.Background(),
		gin.H{"_id": id},
		gin.H{"$set": gin.H{"active": !hidden, "last_activity": util.Now()}},
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
func UpdateJudgeRanking(db *mongo.Database, judge *models.Judge, rankings []primitive.ObjectID) error {
	_, err := db.Collection("judges").UpdateOne(
		context.Background(),
		gin.H{"_id": judge.Id},
		gin.H{"$set": gin.H{"rankings": rankings, "last_activity": util.Now()}},
	)
	return err
}

// TODO: Move the stuff from UpdateJudgeRankings to here
func UpdateJudgeSeenProjects(db *mongo.Database, judge *models.Judge) error {
	_, err := db.Collection("judges").UpdateOne(context.Background(), gin.H{"_id": judge.Id}, gin.H{"$set": gin.H{"seen_projects": judge.SeenProjects}})
	return err
}

// Reset list of projects that judge skipped due to busy
func ResetBusyProjectListForJudge(db *mongo.Database, judge *models.Judge) error {
	_, err := db.Collection("flags").DeleteMany(context.Background(), gin.H{
		"judge_id": judge.Id,
	})
	if err != nil {
		return err
	}
	return nil
}

// Reset project status to non-busy by deleting all "busy" flag for a project
func ResetBusyStatusByProject(db *mongo.Database, project *models.Project) error {
	_, err := db.Collection("flags").DeleteMany(context.Background(), gin.H{
		"project_id": project.Id,
		"reason":     "busy",
	})
	if err != nil {
		return err
	}
	return nil
}

func GetMinJudgeGroup(db *mongo.Database) (int64, error) {
	cursor, err := db.Collection("judges").Aggregate(context.Background(), []gin.H{
		{"$group": gin.H{
			"_id":   "$group",
			"count": gin.H{"$sum": 1},
		}},
	})
	if err != nil {
		return -1, err
	}

	groupCounts := make(map[int64]int64)
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return -1, err
		}
		groupCounts[result["_id"].(int64)] = int64(result["count"].(int32))
	}

	// Sort by count in ascending order
	minGroups := util.SortMapByValue(groupCounts)

	// Get options
	options, err := GetOptions(db, context.Background())
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
func GetNextNJudgeGroups(db *mongo.Database, n int) ([]int64, error) {
	cursor, err := db.Collection("judges").Aggregate(context.Background(), []gin.H{
		{"$group": gin.H{
			"_id":   "$group",
			"count": gin.H{"$sum": 1},
		}},
	})
	if err != nil {
		return nil, err
	}

	groupCounts := make(map[int64]int64)
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		groupCounts[result["_id"].(int64)] = int64(result["count"].(int32))
	}

	// Get options
	options, err := GetOptions(db, context.Background())
	if err != nil {
		return nil, err
	}

	// Find all groups that aren't in groups
	for i := int64(0); i < options.NumGroups; i++ {
		if _, ok := groupCounts[i]; !ok {
			groupCounts[i] = 0
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
