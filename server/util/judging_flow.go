package util

import (
	"math/rand"
	"server/crowdbt"
	"server/database"
	"server/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// Maximum number of tables a judge should visit in a group
const MAX_GROUP_TABLES = 3

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
//  3. Filter projects based on the judge groups
//  4. If there are prioritized projects, pick from that list
//  5. If there are projects not currently being judged, pick from that list
//  6. If there are projects that have less than MIN_VIEWS, pick from that list
//  7. Return the remaining list, filtering at steps 2-4 if applicable
func FindPreferredItems(db *mongo.Database, judge *models.Judge) ([]*models.Project, error) {
	// Get the options from the database
	options, err := database.GetOptions(db)
	if err != nil {
		return nil, err
	}

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

	// Filter projects based on the judge groups
	// Skip this if groups are not enabled or there are no groups
	if options.UseGroups && len(options.Groups) > 0 {
		// First check if the judge has a group and if not assign them a random group
		if len(judge.VisitedGroups) == 0 {
			judge.VisitedGroups = append(judge.VisitedGroups, pickGroup(options, judge))
		}

		// If the judge has visited more than MAX_GROUP_TABLES tables, pick a new group
		if judge.CurrentGroupCount >= MAX_GROUP_TABLES {
			assignJudgeNextGroup(options, judge)
		}

		// Filter out projects that are not in the judge's current group (the last group in the visited groups)
		var groupProjects []*models.Project
		for _, proj := range projects {
			if getProjGroup(options, proj.Location) == judge.VisitedGroups[len(judge.VisitedGroups)-1] {
				groupProjects = append(groupProjects, proj)
			}
		}

		// Increment the judge's current group count
		judge.CurrentGroupCount++

		// If there are projects in the group, pick from that list
		// Otherwise, switch the judge group
		if len(groupProjects) > 0 {
			projects = groupProjects
		} else {
			// TODO: Can we pick from new group bc this currently just ignores the whole grouping thing :(
			assignJudgeNextGroup(options, judge)
		}
	}

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

// assignJudgeNextGroup will assign the judge its next group
func assignJudgeNextGroup(options *models.Options, judge *models.Judge) {
	judge.CurrentGroupCount = 0
	newGroup := pickGroup(options, judge)
	if newGroup == -1 {
		// If there are no more groups, reset the judge's visited groups and pick a first group
		judge.VisitedGroups = []int64{}
		judge.VisitedGroups = append(judge.VisitedGroups, pickGroup(options, judge))
	} else {
		judge.VisitedGroups = append(judge.VisitedGroups, newGroup)
	}
}

// pickGroup picks the next group for the judge to go to.
// Will return -1 if the judge has visited all groups
func pickGroup(op *models.Options, judge *models.Judge) int64 {
	// Filter out groups that the judge has already visited
	var unvisitedGroups []int64
	for groupNum := range op.Groups {
		visited := false
		for _, judgeGroup := range judge.VisitedGroups {
			if int64(groupNum) == judgeGroup {
				visited = true
				break
			}
		}
		if !visited {
			unvisitedGroups = append(unvisitedGroups, int64(groupNum))
		}
	}

	if len(unvisitedGroups) == 0 {
		return -1
	}

	// Randomly pick a group from the unvisited groups
	return unvisitedGroups[rand.Intn(len(unvisitedGroups))]
}

// getProjGroup - Returns the group that the project belongs to.
// Will return -1 if it belongs to no group.
func getProjGroup(op *models.Options, location int64) int64 {
	for groupNum, group := range op.Groups {
		if location >= group.Start && location <= group.End {
			return int64(groupNum)
		}
	}
	return -1
}
