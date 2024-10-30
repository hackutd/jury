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
