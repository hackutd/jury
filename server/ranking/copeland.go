package ranking

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comparison struct {
	Winner primitive.ObjectID `json:"winner"`
	Loser  primitive.ObjectID `json:"loser"`
}

// Convert a judge ranking to a list of pairwise comparisons
func rankingToPairwise(judgeRanking JudgeRanking) []Comparison {
	// Create a slice to store the pairwise comparisons
	pairwise := make([]Comparison, 0)

	// Loop through each project in the ranking and compare it to all the projects below it
	for i, winner := range judgeRanking.Rankings {
		for _, loser := range judgeRanking.Rankings[i+1:] {
			pairwise = append(pairwise, Comparison{winner, loser})
		}
	}

	// Loop through each project in the ranking and compare it to all the unranked projects
	for _, winner := range judgeRanking.Rankings {
		for _, loser := range judgeRanking.Unranked {
			pairwise = append(pairwise, Comparison{winner, loser})
		}
	}

	return pairwise
}

// Calculate the ranking of the projects based on the copeland count method.
// See https://en.wikipedia.org/wiki/Copeland%27s_method
func CalcCopelandRanking(rankingLists []JudgeRanking, projects []primitive.ObjectID) []RankedObject {
	// Create a map to store the scores of each project
	scores := make(map[primitive.ObjectID]float64)

	// Create a list of pairwise comparisons from all judge rankings
	pairs := make([]Comparison, 0)
	for _, rankingList := range rankingLists {
		pairs = append(pairs, rankingToPairwise(rankingList)...)
	}

	// Loop through each pairwise comparison and update the scores of the projects
	for _, pair := range pairs {
		scores[pair.Winner]++
		scores[pair.Loser]--
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
