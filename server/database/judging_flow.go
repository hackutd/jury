package database

import (
	"math/rand"
	"server/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// PickNextProject - Picks the next project for the judge to judge
func PickNextProject(db *mongo.Database, judge *models.Judge, ctx mongo.SessionContext) (*models.Project, error) {
	// Get items
	items, err := FindPreferredItems(db, judge, ctx)
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

	// Pick the first project to be the next project
	return items[0], nil
}

// FindPreferredItems - List of projects to pick from for the judge.
// Find all projects that are higher priority with the following heuristic:
//  1. Ignore all projects that are inactive
//  2. Filter out all projects that the judge has already seen
//  3. Filter out all projects that the judge has flagged (except for busy projects)
//  4. Filter out projects that are currently being judged (if no projects remain after filter, ignore step)
//  5. Filter out all projects that have less than the minimum number of views (if no projects remain after filter, ignore step)
func FindPreferredItems(db *mongo.Database, judge *models.Judge, ctx mongo.SessionContext) ([]*models.Project, error) {
	// Get the list of all active projects
	projects, err := FindActiveProjects(db, ctx)
	if err != nil {
		return nil, err
	}

	// If there are no projects, return an empty list
	if len(projects) == 0 {
		return []*models.Project{}, nil
	}

	// Get all flags for the judge
	flags, err := FindFlagsByJudge(db, judge, ctx)
	if err != nil {
		return nil, err
	}

	// Create a set of voted projects and skipped projects
	done := make(map[string]bool)
	for _, proj := range judge.SeenProjects {
		done[proj.ProjectId.Hex()] = true
	}
	for _, flag := range flags {
		if flag.Reason != "busy" {
			done[flag.ProjectId.Hex()] = true
		}
	}

	// Filter out all projects that the judge has skipped or voted on
	var filteredProjects []*models.Project
	for _, proj := range projects {
		if !done[proj.Id.Hex()] {
			filteredProjects = append(filteredProjects, proj)
		}
	}
	projects = filteredProjects

	// If there are no projects, return an empty list
	// This means that the judge has seen or skipped (except for busy reason) all projects
	if len(projects) == 0 {
		return []*models.Project{}, nil
	}

	// Get all projects currently being judged and remove them from the list
	// Also convert the list to a map for faster lookup
	busyProjects, err := FindBusyProjects(db, ctx)
	if err != nil {
		return nil, err
	}
	busyProjectsMap := make(map[string]bool)
	for _, proj := range busyProjects {
		busyProjectsMap[proj.Hex()] = true
	}

	// Filter out projects that are currently being judged
	// If all projects are busy, ignore this condition
	var freeProjects []*models.Project
	for _, proj := range projects {
		if !busyProjectsMap[proj.Id.Hex()] {
			freeProjects = append(freeProjects, proj)
		}
	}
	if len(freeProjects) > 0 {
		projects = freeProjects
	}

	// Get the minimum number of views of the remaining projects
	minSeen := projects[0].Seen
	for _, proj := range projects {
		if proj.Seen < minSeen {
			minSeen = proj.Seen
		}
	}

	// Filter out projects that have more than the minimum number of views
	var minViewProjects []*models.Project
	for _, proj := range projects {
		if proj.Seen == minSeen {
			minViewProjects = append(minViewProjects, proj)
		}
	}
	if len(minViewProjects) > 0 {
		projects = minViewProjects
	}

	return projects, nil
}
