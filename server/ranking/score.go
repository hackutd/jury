package ranking

import (
	"server/models"
	"slices"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CalculateScoreDiff calculates the score difference between two rankings
// TODO: REMOVE THIS FUNCTION, INSTEAD JUST CALCULATE SCORE AT THE END
func CalculateScoreDiff(rankings []primitive.ObjectID, oldRankings []primitive.ObjectID) *map[primitive.ObjectID]int {
	// Create a map of project IDs to their index in the rankings
	// Use negative numbers here because if there's in the previous but not current then they're removed
	m := len(oldRankings)
	n := len(rankings)

	rankingsMap := make(map[primitive.ObjectID]int)
	for i, proj := range oldRankings {
		rankingsMap[proj] = -(m - i)
	}

	// Iterate through new rankings and add them to the map
	for i, proj := range rankings {
		_, ok := rankingsMap[proj]
		if ok {
			rankingsMap[proj] += n - i
		} else {
			rankingsMap[proj] = n - i
		}
	}

	return &rankingsMap
}

// TODO: Use this for validation of rankings in settings eventually
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
