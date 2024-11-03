package ranking

import (
	"server/models"
	"slices"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CalculateScores(judges []*models.Judge, projects []*models.Project) []RankedObject {
	// Create judge ranking objects
	// Create an array of {Rankings: [], Unranked: []}
	judgeRankings := make([]JudgeRanking, 0)
	for _, judge := range judges {
		unranked := make([]primitive.ObjectID, 0)
		for _, proj := range judge.SeenProjects {
			if !slices.Contains(judge.Rankings, proj.ProjectId) {
				unranked = append(unranked, proj.ProjectId)
			}
		}

		judgeRankings = append(judgeRankings, JudgeRanking{
			Rankings: judge.Rankings,
			Unranked: unranked,
		})
	}

	// Map all projects to their object IDs
	projectIds := make([]primitive.ObjectID, 0)
	for _, proj := range projects {
		projectIds = append(projectIds, proj.Id)
	}

	// Calculate the scores
	return CalcBordaRanking(judgeRankings, projectIds)
}

type ProjectStars struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Stars int                `json:"stars" bson:"stars"`
}

func CalculateStars(judges []*models.Judge, projects []*models.Project) []ProjectStars {
	// Create a list of project stars
	projectStars := make([]ProjectStars, len(projects))
	for i, proj := range projects {
		projectStars[i] = ProjectStars{
			Id:    proj.Id,
			Stars: 0,
		}
	}

	for _, judge := range judges {
		for i, proj := range projects {
			idx := slices.IndexFunc(judge.SeenProjects, func(p models.JudgedProject) bool {
				return p.ProjectId == proj.Id
			})
			if idx != -1 && judge.SeenProjects[idx].Starred {
				projectStars[i].Stars += 1
			}
		}
	}

	// Return the project stars
	return projectStars
}
