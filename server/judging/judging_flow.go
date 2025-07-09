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

func SkipCurrentProject(db *mongo.Database, judge *models.Judge, comps *Comparisons, reason string, getNew bool) error {
	return database.WithTransaction(db, func(sc mongo.SessionContext) error {
		return SkipCurrentProjectWithTx(db, sc, judge, comps, reason, getNew)
	})
}

// SkipCurrentProjectWithTx skips the current project for a judge.
// This is in the judging module instead of the database module to avoid dependency cycles.
// This should be run in a transaction.
func SkipCurrentProjectWithTx(db *mongo.Database, ctx context.Context, judge *models.Judge, comps *Comparisons, reason string, getNew bool) error {
	// Get skipped project from database
	skippedProject, err := database.FindProject(db, ctx, judge.Current)
	if err != nil {
		return errors.New("error finding skipped project in database: " + err.Error())
	}

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

	// Add to flagged project if reason is not busy or absent
	push := make(gin.H)
	if reason != "busy" && reason != "absent" {
		push["flagged"] = skippedProject.Id
	}

	// Update the judge
	_, err = db.Collection("judges").UpdateOne(
		ctx,
		gin.H{"_id": judge.Id},
		gin.H{
			"$set":  gin.H{"current": nil, "last_activity": util.Now()},
			"$push": push,
		},
	)
	if err != nil {
		return err
	}

	// Hide the project if it has been skipped for absent more than 3 times
	err = HideAbsentProject(db, ctx, skippedProject)
	if err != nil {
		return errors.New("error hiding absent project: " + err.Error())
	}

	// Don't get a new project if we're not supposed to
	if !getNew {
		return nil
	}

	// Get a new project
	project, err := PickNextProject(db, ctx, judge, comps)
	if err != nil {
		return err
	}

	if project == nil {
		return nil
	}

	// Update the judge
	return database.UpdateAfterPicked(db, ctx, project, judge)
}

// HideAbsentProject hides a project if it has been absent more than 3 times.
func HideAbsentProject(db *mongo.Database, ctx context.Context, project *models.Project) error {
	projectId := &project.Id

	// Get absent count
	absent, err := database.GetProjectAbsentCount(db, ctx, projectId)
	if err != nil {
		return errors.New("Error getting absent count: " + err.Error())
	}

	// If no more than 3 absences, don't do anything
	if absent < 3 {
		return nil
	}

	// Hide the project
	err = database.SetProjectActive(db, ctx, projectId, false)
	if err != nil {
		return errors.New("Error hiding project: " + err.Error())
	}

	// Delete all absent flags
	err = database.DeleteAbsentFlags(db, ctx, projectId)
	if err != nil {
		return errors.New("Error deleting absent flags: " + err.Error())
	}

	// Add a flag to notate that we are automatically hiding the project
	hideFlag, err := models.NewFlag(project, models.NewDummyJudge(), "hidden-absent")
	if err != nil {
		return errors.New("error creating hidden flag object: " + err.Error())
	}
	err = database.InsertFlag(db, ctx, hideFlag)
	if err != nil {
		return errors.New("error inserting hidden flag into database: " + err.Error())
	}

	return nil
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

	// If the judge has seen more than the proportion of projects in the group, move them to a new group
	if float64(judge.GroupSeen)/float64(numProjects) >= options.AutoSwitchProp {
		judge.Group = (judge.Group + 1) % options.NumGroups
		judge.GroupSeen = 0
	}

	return nil
}

// PickNextProject - Picks the next project for the judge to judge.
// To do this:
//  1. Get all available projects
//  2. If judging a track, simply pick the next project in order
//  3. If any project is prioritized and on the list, return that
//  4. Shuffle projects
//  5. If any projects seen less than min views (set in admin side), only select from that list
//  6. Otherwise, pick the project with the minimum number of comparisons with every other project
func PickNextProject(db *mongo.Database, ctx context.Context, judge *models.Judge, comps *Comparisons) (*models.Project, error) {
	// Get items
	items, err := FindAvailableItems(db, ctx, judge)
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

	// Get options from the db
	options, err := database.GetOptions(db, ctx)
	if err != nil {
		return nil, err
	}

	// If judging a track, simply pick the next project
	if judge.Track != "" {
		return GetNextFreeProject(judge.LastLocation, items)
	}

	// Get prioritized projects
	prioritizedProjects, err := database.GetPrioritizedProjects(db, ctx)
	if err != nil {
		return nil, err
	}

	// If any prioritized projects are in the list, return the first one
	for _, proj := range prioritizedProjects {
		if slices.ContainsFunc(items, func(p *models.Project) bool {
			return p.Id == proj.Id
		}) {
			return proj, nil
		}
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
	if items[0].Seen < options.MinViews {
		return items[0], nil
	}

	// Otherwise, pick the project that has been compared to other projects the least
	return comps.FindLeastCompared(items, judge.SeenProjects), nil
}

// FindAvailableItems - List of projects to pick from for the judge.
// Find all projects that are higher priority with the following heuristic:
//  1. Ignore all projects that are inactive
//  2. Filter out all projects that the judge has already seen
//  3. Filter out all projects that the judge has flagged (except for busy projects)
//  4. Filter out all projects that is not in the judge's track (if tracks are enabled and the user has a track)
//  5. If tracks are enabled, filter out all track projects that have been seen >track_views[track] times
//  6. Filter out projects that are currently being judged (if no projects remain after filter, ignore step)
//  7. If judging a track, return at this point (ignore last 2 conditions)
//  8. Filter out projects not in the judge's group (if no projects remain after filter, try subsequent groups until a project is found OR all projects have been judged)
//  9. Filter out all projects that have less than the minimum number of views (if no projects remain after filter, ignore step)
func FindAvailableItems(db *mongo.Database, ctx context.Context, judge *models.Judge) ([]*models.Project, error) {
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
	flags, err := database.FindFlagsByJudge(db, ctx, judge)
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

	// Get all projects currently being judged
	busyProjects, err := database.FindBusyProjects(db, ctx)
	if err != nil {
		return nil, err
	}

	// Filter out all projects that are not in the judge's track
	// Also filter out all track projects that have been seen >2 times OR is currently being judged by a track judge
	if options.JudgeTracks && judge.Track != "" {
		var trackProjects []*models.Project
		trackIndex := slices.Index(options.Tracks, judge.Track)
		for _, proj := range projects {
			if slices.Contains(proj.ChallengeList, judge.Track) {
				// If currently being judged by the last track judge, do not assign
				if proj.TrackSeen[judge.Track] == options.TrackViews[trackIndex]-1 && busyProjects[proj.Id] == judge.Track {
					continue
				}

				// Otherwise, only add to list if seen < track_view times
				if proj.TrackSeen[judge.Track] < options.TrackViews[trackIndex] {
					trackProjects = append(trackProjects, proj)
				}
			}
		}
		projects = trackProjects
	}

	// If there are no projects, return an empty list
	// This means that the judge has seen or skipped (except for busy reason) all projects
	if len(projects) == 0 {
		return []*models.Project{}, nil
	}

	// Filter out projects that are currently being judged
	// If all projects are busy, ignore this condition
	var freeProjects []*models.Project
	for _, proj := range projects {
		if _, ok := busyProjects[proj.Id]; !ok {
			freeProjects = append(freeProjects, proj)
		}
	}
	if len(freeProjects) > 0 {
		projects = freeProjects
	}

	// If track judging, ignore last conditions
	if judge.Track != "" {
		return projects, nil
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

// GetNextFreeProject - Get the next project in front of the current judges' table number
func GetNextFreeProject(currTableNum int64, projects []*models.Project) (*models.Project, error) {
	// If last location is -1, pick a random project
	if currTableNum == -1 {
		return projects[rand.Intn(len(projects))], nil
	}

	// Sort projects by table number
	slices.SortStableFunc(projects, func(a, b *models.Project) int {
		return int(a.Location - b.Location)
	})

	// Get the next project that is not currently being judged
	for _, proj := range projects {
		if proj.Location > currTableNum {
			return proj, nil
		}
	}

	// This means that the judge is at the last table number; return the first project
	return projects[0], nil
}
