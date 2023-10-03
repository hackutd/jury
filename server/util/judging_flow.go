package util

import (
	"math/rand"
	"server/crowdbt"
	"server/database"
	"server/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// PickNextProject - Picks the next project for the judge to judge
func PickNextProject(db *mongo.Database, judge *models.Judge) (*models.Project, error) {
	// Get items
	items, err := FindPreferredItems(db, judge)
	if err != nil {
		return nil, err
	}

	// If there are no items, return nil
	if len(items) == 0 {
		return nil, nil
	}

	// Shuffle items
	for i := range items {
		j := rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}

	// Randomly pick the first element with an EPSILON chance
	if rand.Float64() < crowdbt.EPSILON {
		return items[0], nil
	}

	// Get previous project
	prevProjectId := judge.Prev
	var prevProject *models.Project
	if prevProjectId == nil {
		prevProject = models.DefaultProject()
	} else {
		prevProject, err = database.FindProjectById(db, prevProjectId)
		if err != nil {
			return nil, err
		}
	}

	// Pick the project with the highest expected information gain
	return maxInfoGain(items, prevProject, judge), nil
}

// FindPreferredItems - List of projects to pick from for the judge.
// Find all projects that are higher priority with the following heuristic:
//  1. Ignore all projects that are inactive
//  2. Filter out all projects that the judge has skipped or voted on
//  3. If there are prioritized projects, pick from that list
//  4. If there are projects not currently being judged, pick from that list
//  5. If there are projects that have less than MIN_VIEWS, pick from that list
//  6. Return the remaining list, filtering at steps 2-4 if applicable
func FindPreferredItems(db *mongo.Database, judge *models.Judge) ([]*models.Project, error) {
	// Get the list of all active projects
	projects, err := database.FindActiveProjects(db)
	if err != nil {
		return nil, err
	}

	// If there are no projects, return an empty list
	if len(projects) == 0 {
		return []*models.Project{}, nil
	}

	// Create a set of skipped/voted projects
	done := make(map[string]bool)
	for _, proj := range judge.SeenProjects {
		done[proj.ProjectId.Hex()] = true
	}

	// Filter out all projects that the judge has skipped or voted on
	var filteredProjects []*models.Project
	for _, proj := range projects {
		if !done[proj.Id.Hex()] {
			filteredProjects = append(filteredProjects, proj)
		}
	}
	projects = filteredProjects

	// If there are prioritized projects, pick from that list
	var prioritizedProjects []*models.Project
	for _, proj := range projects {
		if proj.Prioritized {
			prioritizedProjects = append(prioritizedProjects, proj)
		}
	}
	if len(prioritizedProjects) > 0 {
		projects = prioritizedProjects
	}

	// Get all projects currently being judged and remove them from the list
	// If all projects are busy, ignore this condition
	busyProjects, err := database.FindBusyProjects(db)
	if err != nil {
		return nil, err
	}

	// Convert to map for faster lookup
	busyProjectsMap := make(map[string]bool)
	for _, proj := range busyProjects {
		busyProjectsMap[proj.Hex()] = true
	}

	// Filter out projects that are currently being judged
	var freeProjects []*models.Project
	for _, proj := range projects {
		if !busyProjectsMap[proj.Id.Hex()] {
			freeProjects = append(freeProjects, proj)
		}
	}
	if len(freeProjects) > 0 {
		projects = freeProjects
	}

	// If there are projects that have less than MIN_VIEWS, pick from that list
	var lowViewProjects []*models.Project
	for _, proj := range projects {
		if proj.Seen < crowdbt.MIN_VIEWS {
			lowViewProjects = append(lowViewProjects, proj)
		}
	}
	if len(lowViewProjects) > 0 {
		projects = lowViewProjects
	}

	return projects, nil
}

// maxInfoGain - Returns the project with the maximum information gain
func maxInfoGain(items []*models.Project, prev *models.Project, judge *models.Judge) *models.Project {
	maxInfoGain := -1.0
	var maxInfoGainProject *models.Project
	for _, item := range items {
		infoGain := crowdbt.ExpectedInformationGain(judge.Alpha, judge.Beta, prev.Mu, prev.SigmaSq, item.Mu, item.SigmaSq)
		if infoGain > maxInfoGain {
			maxInfoGain = infoGain
			maxInfoGainProject = item
		}
	}
	return maxInfoGainProject
}
