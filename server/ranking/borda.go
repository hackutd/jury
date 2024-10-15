package ranking

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CalcBordaRanking calculates the ranking of projects based on the borda count rank aggregation model.
// See https://en.wikipedia.org/wiki/Borda_count
func CalcBordaRanking(rankingLists []JudgeRanking, projects []primitive.ObjectID) []RankedObject {
	// Create a map to store the scores of each project
	scores := make(map[primitive.ObjectID]float64)

	// Loop through judges' rankings
	for _, rankingList := range rankingLists {
		// Give n points to 1st place, n-1 to 2nd place, ... 1 to last place,
		// where n is the number of ranked projects
		n := len(rankingList.Rankings)

		for i, projId := range rankingList.Rankings {
			scores[projId] += float64(n - i)
		}

	}

	// Create the output DS
	ranked := make([]RankedObject, 0)
	for _, project := range projects {
		ranked = append(ranked, RankedObject{project, scores[project]})
	}

	// Sort the projects by their scores
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}
