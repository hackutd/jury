package judging

import (
	"context"
	"server/database"
	"server/models"
	"slices"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AggregateRanking will take the ranking array from the judge and create
// the aggregated ranking model using the Copeland counting method
func AggregateRanking(judge *models.Judge) []models.AggRanking {
	// Get unranked projects
	var unranked []primitive.ObjectID
	for _, p := range judge.SeenProjects {
		if !slices.Contains(judge.Rankings, p.ProjectId) {
			unranked = append(unranked, p.ProjectId)
		}
	}

	// Create output
	rc := len(judge.Rankings)
	uc := len(unranked)
	out := make([]models.AggRanking, 0, rc+uc)

	// Calculate scores for ranked projects
	// score = won - loss + unranked count
	// score = (rc - i - 1) - i + uc = rc + uc - 1 - 2i
	// c = rc + uc - 1; score = c - 2i
	c := rc + uc - 1
	for i, p := range judge.Rankings {
		out = append(out, *models.NewAggRanking(p, int64(c-2*i)))
	}

	// Calculate scores for unranked projects
	// score = -rc (beaten by all ranked projects)
	for _, p := range unranked {
		out = append(out, *models.NewAggRanking(p, int64(-rc)))
	}

	return out
}

// InitAggregateRankings will re-calculate all judge ranking aggregation objects
// when the Jury backend starts
func InitAggregateRankings(db *mongo.Database) error {
	return database.WithTransaction(db, func(sc mongo.SessionContext) error {
		// Get all judges
		judges, err := database.FindAllJudges(db, sc)
		if err != nil {
			return err
		}

		// Use bulk write to update all judges' aggregated rankings
		models := make([]mongo.WriteModel, 0, len(judges))
		for _, judge := range judges {
			agg := AggregateRanking(judge)
			models = append(models, mongo.NewUpdateOneModel().SetFilter(gin.H{"_id": judge.Id}).SetUpdate(gin.H{"$set": gin.H{"rankings_agg": agg}}))
		}
		opts := options.BulkWrite().SetOrdered(false)
		_, err = db.Collection("judges").BulkWrite(sc, models, opts)
		return err
	})
}

type ProjectScores struct {
	Score      int64            `bson:"score" json:"score"`
	Stars      int64            `bson:"stars" json:"stars"`
	TrackStars map[string]int64 `bson:"track_stars" json:"track_stars"`
}

type ProjectScoresWithId struct {
	ProjectId  primitive.ObjectID `bson:"_id" json:"_id"`
	Score      int64              `bson:"score" json:"score"`
	Stars      int64              `bson:"stars" json:"stars"`
	TrackStars map[string]int64   `bson:"track_stars" json:"track_stars"`
}

func removeId(scoresWithId *ProjectScoresWithId) *ProjectScores {
	return &ProjectScores{
		Score:      scoresWithId.Score,
		Stars:      scoresWithId.Stars,
		TrackStars: scoresWithId.TrackStars,
	}
}

// AggregateScores takes the scores and stars from judges and aggregates them to form a final ranking
func AggregateScores(db *mongo.Database, ctx context.Context) (map[primitive.ObjectID]ProjectScores, error) {
	pipeline := mongo.Pipeline{
		// === Pipeline 1: Aggregate all general judges' ranking scores ===
		bson.D{{"$match", bson.D{{"track", ""}}}},

		bson.D{{"$unwind", "$rankings_agg"}},
		bson.D{{"$group", bson.D{
			{"_id", "$rankings_agg.project_id"},
			{"score", bson.D{{"$sum", "$rankings_agg.score"}}},
		}}},

		// === Pipeline 2: Aggregate all general judges' stars ===
		bson.D{{"$unionWith", gin.H{
			"coll": "judges",
			"pipeline": []gin.H{
				{"$match": gin.H{"track": ""}},
				{"$unwind": "$seen_projects"},
				{"$match": gin.H{"seen_projects.starred": true}},
				{"$group": gin.H{
					"_id":   "$seen_projects.project_id",
					"stars": gin.H{"$sum": 1},
				}},
			},
		}}},

		// === Combine scores and stars (first stage) ===
		bson.D{{"$group", bson.D{
			{"_id", "$_id"},
			{"score", bson.D{{"$sum", "$score"}}},
			{"stars", bson.D{{"$sum", "$stars"}}},
		}}},

		// // === Pipeline 3: Aggregate all track judges' stars, grouped by project and track ===
		bson.D{{"$unionWith", gin.H{
			"coll": "judges",
			"pipeline": []gin.H{
				{"$match": gin.H{"track": gin.H{"$ne": ""}}},
				{"$unwind": "$seen_projects"},
				{"$match": gin.H{"seen_projects.starred": true}},
				{"$group": gin.H{
					"_id": gin.H{
						"projectId": "$seen_projects.project_id",
						"track":     "$track",
					},
					"count": gin.H{"$sum": 1},
				}},
				{"$group": gin.H{
					"_id": "$_id.projectId",
					"track_stars": gin.H{
						"$push": gin.H{
							"k": "$_id.track",
							"v": "$count",
						},
					},
				}},
				{"$addFields": gin.H{
					"track_stars": gin.H{
						"$arrayToObject": "$track_stars",
					},
				}},
			},
		}}},

		// // === Final merge of all fields ===
		bson.D{{"$group", bson.D{
			{"_id", "$_id"},
			{"score", bson.D{{"$sum", "$score"}}},
			{"stars", bson.D{{"$sum", "$stars"}}},
			{"track_stars", bson.D{{"$mergeObjects", "$track_stars"}}},
		}}},
	}

	// Run the aggregation and get results with cursor
	cursor, err := db.Collection("judges").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []ProjectScoresWithId
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Convert results to map
	out := make(map[primitive.ObjectID]ProjectScores)
	for _, p := range results {
		out[p.ProjectId] = *removeId(&p)
	}

	return out, nil
}
