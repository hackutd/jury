package database

import (
	"context"
	"errors"
	"server/models"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateProjectLastActivity to the current time
func UpdateProjectLastActivity(db *mongo.Database, ctx context.Context, id *primitive.ObjectID) error {
	// Get current time
	lastActivity := util.Now()
	_, err := db.Collection("projects").UpdateOne(ctx, gin.H{"_id": id}, gin.H{"$set": gin.H{"last_activity": lastActivity}})
	return err
}

// InsertProjects inserts a list of projects into the database
func InsertProjects(db *mongo.Database, projects []*models.Project) error {
	var docs []interface{}
	for _, project := range projects {
		docs = append(docs, project)
	}
	_, err := db.Collection("projects").InsertMany(context.Background(), docs)
	return err
}

// InsertProject inserts a project into the database
func InsertProject(db *mongo.Database, ctx context.Context, project *models.Project) error {
	_, err := db.Collection("projects").InsertOne(ctx, project)
	return err
}

// FindAllProjects returns a list of all projects in the database
func FindAllProjects(db *mongo.Database, ctx context.Context) ([]*models.Project, error) {
	projects := make([]*models.Project, 0)
	cursor, err := db.Collection("projects").Find(ctx, gin.H{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// FindProjectsByGroup returns a list of all projects with the specific track
func FindProjectsByTrack(db *mongo.Database, ctx context.Context, track string) ([]*models.Project, error) {
	projects := make([]*models.Project, 0)
	cursor, err := db.Collection("projects").Find(ctx, gin.H{"challenge_list": track})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetPrioritizedProjects returns a list of all prioritized projects in the database
func GetPrioritizedProjects(db *mongo.Database, ctx context.Context) ([]*models.Project, error) {
	projects := make([]*models.Project, 0)
	cursor, err := db.Collection("projects").Find(ctx, gin.H{"prioritized": true}, options.Find().SetProjection(gin.H{"_id": 1}))
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// DeleteProjectById deletes a project from the database by id
func DeleteProjectById(db *mongo.Database, id primitive.ObjectID) error {
	_, err := db.Collection("projects").DeleteOne(context.Background(), gin.H{"_id": id})
	return err
}

// AggregateProjectStats aggregates all stats from the database for a project
func AggregateProjectStats(db *mongo.Database) (*models.ProjectStats, error) {
	// Get the totoal number of projects
	totalProjects, err := db.Collection("projects").EstimatedDocumentCount(context.Background())
	if err != nil {
		return nil, err
	}

	// Get the average votes and average seen using an aggregation pipeline
	cursor, err := db.Collection("projects").Aggregate(context.Background(), []gin.H{
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
	var stats models.ProjectStats
	cursor.Next(context.Background())
	err = cursor.Decode(&stats)
	if err != nil {
		if err.Error() == "EOF" {
			stats = models.ProjectStats{Num: 0, AvgSeen: 0, NumActive: 0}
		} else {
			return nil, err
		}
	}

	// Set the total number of projects
	stats.Num = totalProjects

	return &stats, nil
}

// FindActiveProjects returns a list of all active projects in the database
func FindActiveProjects(db *mongo.Database, ctx mongo.SessionContext) ([]*models.Project, error) {
	var projects []*models.Project
	cursor, err := db.Collection("projects").Find(ctx, gin.H{"active": true})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// FindBusyProjects returns a list of all projects that are currently being judged.
// To do this, we collect all projects in the judge's "current" field.
// We will return a map of project IDs to track names.
func FindBusyProjects(db *mongo.Database, ctx mongo.SessionContext) (map[primitive.ObjectID]string, error) {
	// Get all judges that are currently judging a project
	var judges []*models.Judge
	cursor, err := db.Collection("judges").Find(ctx, gin.H{
		"current": gin.H{
			"$ne": nil,
		},
		"active": true,
	}, options.Find().SetProjection(gin.H{"current": 1, "track": 1}))
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &judges)
	if err != nil {
		return nil, err
	}

	// Extract the project IDs from the judges
	output := make(map[primitive.ObjectID]string)
	for _, judge := range judges {
		if judge.Current != nil {
			output[*judge.Current] = judge.Track
		}
	}
	return output, nil
}

// FindProjectById returns a project from the database by id
func FindProjectById(db *mongo.Database, ctx context.Context, id *primitive.ObjectID) (*models.Project, error) {
	var project models.Project
	err := db.Collection("projects").FindOne(ctx, gin.H{"_id": id}).Decode(&project)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateAfterPicked updates the seen value of the new project picked and the judge's current project
func UpdateAfterPicked(db *mongo.Database, project *models.Project, judge *models.Judge) error {
	err := WithTransaction(db, func(ctx mongo.SessionContext) error {
		return UpdateAfterPickedWithTx(db, ctx, project, judge)
	})
	return err
}

// UpdateAfterPickedWithTx updates the seen value of the new project picked and the judge's current project
func UpdateAfterPickedWithTx(db *mongo.Database, ctx context.Context, project *models.Project, judge *models.Judge) error {
	// De-prioritize project
	_, err := db.Collection("projects").UpdateOne(
		ctx,
		gin.H{"_id": project.Id},
		gin.H{"$set": gin.H{"prioritized": false}},
	)
	if err != nil {
		return err
	}

	// Set the judge's current project
	_, err = db.Collection("judges").UpdateOne(
		ctx,
		gin.H{"_id": judge.Id},
		gin.H{"$set": gin.H{"last_location": project.Location, "current": project.Id, "last_activity": util.Now()}},
	)
	return err
}

// CountProjectDocuments returns the number of documents in the projects collection
func CountProjectDocuments(db *mongo.Database) (int64, error) {
	return db.Collection("projects").EstimatedDocumentCount(context.Background())
}

// CountTrackProjects returns the number of projects in a specific track
func CountTrackProjects(db *mongo.Database, track string) (int64, error) {
	return db.Collection("projects").CountDocuments(context.Background(), gin.H{"challenge_list": track})
}

// SetProjectActive sets the active field of a project (hide or unhide project)
func SetProjectActive(db *mongo.Database, ctx context.Context, id *primitive.ObjectID, active bool) error {
	_, err := db.Collection("projects").UpdateOne(context.Background(), gin.H{"_id": id}, gin.H{"$set": gin.H{"active": active}})
	return err
}

// SetProjectPrioritized sets the prioritized field of a project
func SetProjectPrioritized(db *mongo.Database, id *primitive.ObjectID, prioritized bool) error {
	_, err := db.Collection("projects").UpdateOne(context.Background(), gin.H{"_id": id}, gin.H{"$set": gin.H{"prioritized": prioritized}})
	return err
}

// UpdateProjects will update ALL projects in the database
func UpdateProjects(db *mongo.Database, projects []*models.Project) error {
	models := make([]mongo.WriteModel, 0, len(projects))
	for _, project := range projects {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(gin.H{"_id": project.Id}).SetUpdate(gin.H{"$set": project}))
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := db.Collection("projects").BulkWrite(context.Background(), models, opts)
	return err
}

// DecrementProjectSeenCount decrements the seen count of a project (after being skipped)
func DecrementProjectSeenCount(db *mongo.Database, ctx context.Context, project *models.Project) error {
	_, err := db.Collection("projects").UpdateOne(ctx, gin.H{"_id": project.Id}, gin.H{"$inc": gin.H{"seen": -1}})
	return err
}

// GetNumProjectsInGroup returns the number of projects in a group
// TODO: Can we pre-aggregate this value?
func GetNumProjectsInGroup(db *mongo.Database, ctx context.Context, group int64) (int64, error) {
	return db.Collection("projects").CountDocuments(ctx, gin.H{"group": group})
}

// GetChallenges gets the list of all challenges from the database
func GetChallenges(db *mongo.Database, ctx context.Context) ([]string, error) {
	// Get all projects
	projects, err := FindAllProjects(db, ctx)
	if err != nil {
		return nil, err
	}

	// Extract the challenges from the projects
	challenges := make(map[string]bool)
	for _, project := range projects {
		for _, challenge := range project.ChallengeList {
			challenges[challenge] = true
		}
	}

	// Convert the map to a list
	var challengeList []string
	for challenge := range challenges {
		challengeList = append(challengeList, challenge)
	}

	return challengeList, nil
}

// UpdateProjectStars updates the starred count of a project, incrementing or decrementing it
func UpdateProjectStars(db *mongo.Database, ctx context.Context, projId primitive.ObjectID, increment bool) error {
	change := -1
	if increment {
		change = 1
	}

	_, err := db.Collection("projects").UpdateOne(ctx, gin.H{"_id": projId}, gin.H{"$inc": gin.H{"stars": change}})
	return err
}

// UpdateProjectTrackStars updates the starred count of a project in a specific track, incrementing or decrementing it
func UpdateProjectTrackStars(db *mongo.Database, ctx context.Context, projId primitive.ObjectID, track string, increment bool) error {
	change := -1
	if increment {
		change = 1
	}
	track_str := "track_stars." + track

	_, err := db.Collection("projects").UpdateOne(ctx, gin.H{"_id": projId}, gin.H{"$inc": gin.H{track_str: change}})
	return err
}

// UpdateProjectScores updates the scores of the projects given the diff
func UpdateProjectScores(db *mongo.Database, ctx context.Context, scoreDiff *map[primitive.ObjectID]int) error {
	models := make([]mongo.WriteModel, 0, len(*scoreDiff))
	for id, diff := range *scoreDiff {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(gin.H{"_id": id}).SetUpdate(gin.H{"$inc": gin.H{"score": diff}}))
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := db.Collection("projects").BulkWrite(ctx, models, opts)
	if err != nil {
		return errors.New("error updating projects in database: " + err.Error())
	}

	return nil
}

// GetMaxTableNum returns the max table number assigned to a project;
// This will be the maximum table number in the database
func GetMaxTableNum(db *mongo.Database, ctx context.Context) (int64, error) {
	// Get the maximum table number
	cursor, err := db.Collection("projects").Find(ctx, gin.H{}, options.Find().SetSort(gin.H{"location": -1}).SetLimit(1))
	if err != nil {
		return 0, err
	}

	var project models.Project
	if cursor.Next(ctx) {
		err = cursor.Decode(&project)
		if err != nil {
			return 0, err
		}
	}

	return project.Location, nil
}

func ReassignAllGroupNums(db *mongo.Database, ctx context.Context, op *models.Options) error {
	// If multi-group is not enabled, set all projects to group 0
	if !op.MultiGroup {
		return nil
	}

	// Get all projects
	projects, err := FindAllProjects(db, ctx)
	if err != nil {
		return err
	}

	// Create bulk write
	models := make([]mongo.WriteModel, 0, len(projects))

	// Reassign group numbers to all projects
	for _, project := range projects {
		group := util.GroupFromTable(op, project.Location)
		models = append(models, mongo.NewUpdateOneModel().SetFilter(gin.H{"_id": project.Id}).SetUpdate(gin.H{"$set": gin.H{"group": group}}))
	}

	// Bulk write
	opts := options.BulkWrite().SetOrdered(false)
	_, err = db.Collection("projects").BulkWrite(ctx, models, opts)
	if err != nil {
		return errors.New("error updating projects in database: " + err.Error())
	}

	return nil
}
