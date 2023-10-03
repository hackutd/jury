package database

import (
	"context"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

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
func InsertProject(db *mongo.Database, project *models.Project) error {
	_, err := db.Collection("projects").InsertOne(context.Background(), project)
	return err
}

// FindAllProjects returns a list of all projects in the database
func FindAllProjects(db *mongo.Database) ([]*models.Project, error) {
	projects := make([]*models.Project, 0)
	cursor, err := db.Collection("projects").Find(context.Background(), gin.H{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &projects)
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
			"avgVotes": gin.H{
				"$avg": "$votes",
			},
			"avgSeen": gin.H{
				"$avg": "$seen",
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
		return nil, err
	}

	// Set the total number of projects
	stats.Num = totalProjects

	return &stats, nil
}

// FindActiveProjects returns a list of all active projects in the database
func FindActiveProjects(db *mongo.Database) ([]*models.Project, error) {
	var projects []*models.Project
	cursor, err := db.Collection("projects").Find(context.Background(), gin.H{"active": true})
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
// To do this, we collect all projects in the judge's "next" field
func FindBusyProjects(db *mongo.Database) ([]*primitive.ObjectID, error) {
	var judges []*models.Judge
	cursor, err := db.Collection("judges").Find(context.Background(), gin.H{
		"next": gin.H{
			"$ne": nil,
		},
		"active": true,
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &judges)
	if err != nil {
		return nil, err
	}

	// Extract the project IDs from the judges
	var projects []*primitive.ObjectID
	for _, judge := range judges {
		projects = append(projects, judge.Next)
	}
	return projects, nil
}

// FindProjectById returns a project from the database by id
func FindProjectById(db *mongo.Database, id *primitive.ObjectID) (*models.Project, error) {
	var project models.Project
	err := db.Collection("projects").FindOne(context.Background(), gin.H{"_id": id}).Decode(&project)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// TODO: Abstract transaction to a function wrapper

// Update the seen value of the new project picked and the seen list of the judge that saw it
func UpdateProjectSeen(db *mongo.Database, project *models.Project, judge *models.Judge) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		// Update the project's seen value
		_, err := db.Collection("projects").UpdateOne(context.Background(), gin.H{"_id": project.Id}, gin.H{"$inc": gin.H{"seen": 1}})
		if err != nil {
			return nil, err
		}

		// Add the project to the judge's seen list
		judge.SeenProjects = append(judge.SeenProjects, *models.JudgeProjectFromProject(project))

		// Update the judge
		_, err = db.Collection("judges").UpdateOne(context.Background(), gin.H{"_id": judge.Id}, gin.H{"$set": gin.H{"seen_projects": judge.SeenProjects}})

		return nil, err
	}, txnOptions)

	return err
}

// CountProjectDocuments returns the number of documents in the projects collection
func CountProjectDocuments(db *mongo.Database) (int64, error) {
	return db.Collection("projects").EstimatedDocumentCount(context.Background())
}

// SetProjectHidden sets the active field of a judge
func SetProjectHidden(db *mongo.Database, id *primitive.ObjectID, hidden bool) error {
	_, err := db.Collection("projects").UpdateOne(context.Background(), gin.H{"_id": id}, gin.H{"$set": gin.H{"active": !hidden}})
	return err
}

// SetProjectPrioritized sets the prioritized field of a project
func SetProjectPrioritized(db *mongo.Database, id *primitive.ObjectID, prioritized bool) error {
	_, err := db.Collection("projects").UpdateOne(context.Background(), gin.H{"_id": id}, gin.H{"$set": gin.H{"prioritized": prioritized}})
	return err
}
