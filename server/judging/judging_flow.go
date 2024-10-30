package judging

import (
	"context"
	"errors"
	"math/rand"
	"server/database"
	"server/models"
	"server/util"
	"slices"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SkipCurrentProject skips the current project for a judge.
// This is in the judging module instead of the database module to avoid dependency cycles.
func SkipCurrentProject(db *mongo.Database, judge *models.Judge, comps *Comparisons, reason string, getNew bool) error {
	return database.WithTransaction(db, func(ctx mongo.SessionContext) error {
		// Get skipped project from database
		skippedProject, err := database.FindProjectById(db, ctx, judge.Current)
		if err != nil {
			return errors.New("error finding skipped project in database: " + err.Error())
		}

		// If skipping for any reason other than wanting a break, add the project to the skipped list
		if reason != "break" {
			// Create a new skip object
			skip, err := models.NewFlag(skippedProject, judge, reason)
			if err != nil {
				return errors.New("error creating flag object: " + err.Error())
			}

			// Add skipped project to flags database
			err = database.InsertFlag(db, ctx, skip)
			if err != nil {
				return errors.New("error inserting flag into database: " + err.Error())
			}
		}

		// Update the judge
		_, err = db.Collection("judges").UpdateOne(
			ctx,
			gin.H{"_id": judge.Id},
			gin.H{"$set": gin.H{"current": nil, "last_activity": util.Now()}},
		)
		if err != nil {
			return err
		}

		// Check to see how many times the project has been marked absent
		active := true
		absentCount, err := database.GetProjectAbsentCount(db, skippedProject, ctx)
		if err != nil {
			return errors.New("error checking for project flags: " + err.Error())
		}
		// Hide project if it's been absent more than 3 times
		if reason == "absent" && absentCount >= 2 {
			active = false

			// Also remove all absent flags for the project
			err = database.DeleteAbsentFlags(db, skippedProject, ctx)
			if err != nil {
				return errors.New("error deleting absent flags: " + err.Error())
			}
		}

		// Update the project
		_, err = db.Collection("projects").UpdateOne(ctx, gin.H{"_id": skippedProject.Id}, gin.H{"$inc": gin.H{"seen": -1}, "$set": gin.H{"active": active}})
		if err != nil {
			return err
		}

		// Don't get a new project if we're not supposed to
		if !getNew {
			return nil
		}

		// Get a new project
		project, err := PickNextProject(db, judge, ctx, comps)
		if err != nil {
			return err
		}

		if project == nil {
			return nil
		}

		// Update the judge
		return database.UpdateAfterPickedWithTx(db, project, judge, ctx)
	})
}

// MoveJudgeGroup will increment the count of projects they've seen in the current group,
// then moves a judge to a new group if they've seen n projects in their current group,
// where n is either the count or the proportion of projects in the group.
func MoveJudgeGroup(db *mongo.Database, ctx context.Context, judge *models.Judge, options *models.Options) error {
	// Get the number of projects in the group
	numProjects, err := database.GetNumProjectsInGroup(db, ctx, judge.Group)
	if err != nil {
		return errors.New("error getting number of projects in group: " + err.Error())
	}

	// Increment the count of projects seen in the group
	judge.GroupSeen++

	// Make sure the group switch method is valid
	if options.MainGroup.AutoSwitchMethod != "count" && options.MainGroup.AutoSwitchMethod != "proportion" {
		return errors.New("invalid group switch method detected: " + options.MainGroup.AutoSwitchMethod)
	}

	// If the judge has seen more than the number of projects in the group, move them to a new group
	if (options.MainGroup.AutoSwitchMethod == "count" && judge.GroupSeen >= options.MainGroup.AutoSwitchCount) ||
		(options.MainGroup.AutoSwitchMethod == "proportion" && float64(judge.GroupSeen)/float64(numProjects) >= options.MainGroup.AutoSwitchProp) {
		judge.Group = (judge.Group + 1) % options.NumGroups
		judge.GroupSeen = 0
	}

	return nil
}

// PickNextProject - Picks the next project for the judge to judge.
// To do this:
//  1. Shuffle projects
//  2. If any projects seen less than min views (set in admin side), only select from that list
//  3. Otherwise, pick the project with the minimum number of comparisons with every other project
func PickNextProject(db *mongo.Database, judge *models.Judge, ctx mongo.SessionContext, comps *Comparisons) (*models.Project, error) {
	// Get items
	items, err := FindPreferredItems(db, judge, ctx)
	if err != nil {
		return nil, err
	}

	// If there are no items, return nil
	if len(items) == 0 {
		return nil, nil
	}

	// If there is only one item, return that
	if len(items) == 1 {
		return items[0], nil
	}

	// Get min views from db
	minViews, err := database.GetMinViews(db)
	if err != nil {
		return nil, err
	}

	// Shuffle items
	for i := range items {
		j := rand.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}

	// Stable sort by the number of views
	slices.SortStableFunc(items, func(a, b *models.Project) int {
		return int(a.Seen - b.Seen)
	})

	// If any items have not been seen minViews times, return that
	// This will be a random item due to shuffling + stable sort
	if items[0].Seen < minViews {
		return items[0], nil
	}

	// Otherwise, pick the project that has been compared to other projects the least
	return comps.FindLeastCompared(items, judge.SeenProjects), nil
}

// FindPreferredItems - List of projects to pick from for the judge.
// Find all projects that are higher priority with the following heuristic:
//  1. Ignore all projects that are inactive
//  2. Filter out all projects that the judge has already seen
//  3. Filter out all projects that the judge has flagged (except for busy projects)
//  4. Filter out projects that are currently being judged (if no projects remain after filter, ignore step)
//  5. Filter out projects not in the judge's group (if no projects remain after filter, try subsequent groups until a project is found OR all projects have been judged)
//  6. Filter out all projects that have less than the minimum number of views (if no projects remain after filter, ignore step)
func FindPreferredItems(db *mongo.Database, judge *models.Judge, ctx mongo.SessionContext) ([]*models.Project, error) {
	// Get the list of all active projects
	projects, err := database.FindActiveProjects(db, ctx)
	if err != nil {
		return nil, err
	}

	// If there are no projects, return an empty list
	if len(projects) == 0 {
		return []*models.Project{}, nil
	}

	// Get all flags for the judge
	flags, err := database.FindFlagsByJudge(db, judge, ctx)
	if err != nil {
		return nil, err
	}

	// Get the options
	options, err := database.GetOptions(db, ctx)
	if err != nil {
		return nil, err
	}

	// Create a set of voted projects and skipped projects
	done := make(map[string]bool)
	for _, proj := range judge.SeenProjects {
		done[proj.ProjectId.Hex()] = true
	}
	for _, flag := range flags {
		done[flag.ProjectId.Hex()] = true
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
	busyProjects, err := database.FindBusyProjects(db, ctx)
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

	// Group checking logic
	if options.MultiGroup {
		attempts := int64(0)
		for {
			// Filter out projects not in judge's group
			var groupProjects []*models.Project
			for _, proj := range projects {
				if proj.Group == judge.Group {
					groupProjects = append(groupProjects, proj)
				}
			}

			// If projects exitst, return them
			if len(groupProjects) > 0 {
				projects = groupProjects
				break
			}

			// If no projects exist, move to the next group
			judge.Group = (judge.Group + 1) % options.NumGroups
			attempts++

			// If we've tried all groups, break
			if attempts >= options.NumGroups {
				break
			}
		}
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
