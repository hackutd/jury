package router

import (
	"errors"
	"net/http"
	"server/database"
	"server/funcs"
	"server/judging"
	"server/models"
	"server/util"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

// POST /project/devpost - AddDevpostCsv adds a csv export from devpost to the database
func AddDevpostCsv(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the CSV file from the request
	file, err := ctx.FormFile("csv")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading CSV file from request: " + err.Error()})
		return
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error opening CSV file: " + err.Error()})
		return
	}

	// Read the file
	content := make([]byte, file.Size)
	_, err = f.Read(content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reading CSV file: " + err.Error()})
		return
	}

	// Parse the CSV file
	projects, err := funcs.ParseDevpostCSV(string(content), state.Db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing CSV file: " + err.Error()})
		return
	}

	// Run the remaining actions in a transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Insert projects into the database
		err = database.InsertProjects(state.Db, sc, projects)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting projects into database: " + err.Error()})
			return err
		}

		// Reload comparisons
		err = judging.ReloadComparisons(state.Db, sc, state.Comps)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reloading comparisons: " + err.Error()})
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	// Send OK
	state.Logger.AdminLogf("Added %d projects from Devpost CSV", len(projects))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /project/new - AddProject adds a project to the database
func AddProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the projectReq from the request
	var projectReq models.AddProjectRequest
	err := ctx.BindJSON(&projectReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error binding project from request: " + err.Error()})
		return
	}

	// Make sure name, description, and url are defined
	if projectReq.Name == "" || projectReq.Description == "" || projectReq.Url == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name, description, and url are required"})
		return
	}

	// Get the challenge list
	challengeList := strings.Split(projectReq.ChallengeList, ",")
	if projectReq.ChallengeList == "" {
		challengeList = []string{}
	}
	for i := range challengeList {
		challengeList[i] = strings.TrimSpace(challengeList[i])
	}

	var project *models.Project

	// Run the remaining actions in a transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Get the options from the database
		options, err := database.GetOptions(state.Db, sc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options from database: " + err.Error()})
			return err
		}

		// If the challenge list contains the ignore track, skip the project
		ignore := false
		for _, ignoreTrack := range options.IgnoreTracks {
			if slices.Contains(challengeList, ignoreTrack) {
				ignore = true
				break
			}
		}
		if ignore {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "project is in an ignored track"})
			return errors.New("project is in an ignored track")
		}

		// Get max project number
		tableNum, err := database.GetMaxTableNum(state.Db, sc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting max project number from database: " + err.Error()})
			return err
		}

		// Create the project
		project = models.NewProject(projectReq.Name, tableNum+1, util.GroupFromTable(options, tableNum+1), projectReq.Description, projectReq.Url, projectReq.TryLink, projectReq.VideoLink, challengeList)

		// Insert project
		err = database.InsertProject(state.Db, sc, project)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting project into database: " + err.Error()})
			return err
		}

		// Reload comparisons
		err = judging.ReloadComparisons(state.Db, sc, state.Comps)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reloading comparisons: " + err.Error()})
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Send OK
	state.Logger.AdminLogf("Added project %s (%s)", project.Name, project.Id.Hex())
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /project/list - ListProjects lists all projects in the database
func ListProjects(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the projects from the database
	projects, err := database.FindAllProjects(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects from database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, projects)
}

type PublicProject struct {
	Name          string `json:"name"`
	Location      int64  `json:"location"`
	Group         int64  `json:"group"`
	Description   string `json:"description"`
	Url           string `json:"url"`
	TryLink       string `json:"try_link"`
	VideoLink     string `json:"video_link"`
	ChallengeList string `json:"challenge_list"`
}

func ListPublicProjects(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the projects from the database
	projects, err := database.FindAllProjects(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects from database: " + err.Error()})
		return
	}

	// Convert projects to public projects
	publicProjects := make([]PublicProject, len(projects))
	for i, project := range projects {
		publicProjects[i] = PublicProject{
			Name:          project.Name,
			Location:      project.Location,
			Group:         project.Group,
			Description:   project.Description,
			Url:           project.Url,
			TryLink:       project.TryLink,
			VideoLink:     project.VideoLink,
			ChallengeList: strings.Join(project.ChallengeList, ", "),
		}
	}

	// Send OK
	ctx.JSON(http.StatusOK, publicProjects)
}

// POST /project/csv - Endpoint to add projects from a CSV file
func AddProjectsCsv(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the CSV file from the request
	file, err := ctx.FormFile("csv")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading CSV file from request: " + err.Error()})
		return
	}

	// Get the hasHeader parameter from the request
	hasHeader := ctx.PostForm("hasHeader") == "true"

	// Open the file
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error opening CSV file: " + err.Error()})
		return
	}

	// Read the file
	content := make([]byte, file.Size)
	_, err = f.Read(content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reading CSV file: " + err.Error()})
		return
	}

	// Parse the CSV file
	projects, err := funcs.ParseProjectCsv(string(content), hasHeader, state.Db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing CSV file: " + err.Error()})
		return
	}

	// Run the remaining actions in a transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Insert projects into the database
		err = database.InsertProjects(state.Db, sc, projects)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting projects into database: " + err.Error()})
			return err
		}

		// Reload comparisons
		err = judging.ReloadComparisons(state.Db, sc, state.Comps)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reloading comparisons: " + err.Error()})
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Send OK
	state.Logger.AdminLogf("Added %d projects from CSV", len(projects))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// PUT /project/:id - EditProject edits a project in the database
func EditProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the body content
	var projReq models.AddProjectRequest
	err := ctx.BindJSON(&projReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the challenge list
	challengeList := strings.Split(projReq.ChallengeList, ",")
	if projReq.ChallengeList == "" {
		challengeList = []string{}
	}
	for i := range challengeList {
		challengeList[i] = strings.TrimSpace(challengeList[i])
	}

	// Get the id from the request
	id := ctx.Param("id")

	// Convert project ID string to ObjectID
	projectObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Edit the project in the database
	err = database.UpdateProjectBasicInfo(state.Db, projectObjectId, &projReq, challengeList)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project in database: " + err.Error()})
		return
	}

	state.Logger.AdminLogf("Edited project %s: %s", id, projReq.Name)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// DELETE /project/:id - DeleteProject deletes a project from the database
func DeleteProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the id from the request
	id := ctx.Param("id")

	// Convert judge ID string to ObjectID
	projectObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Delete the project from the database
	err = database.DeleteProjectById(state.Db, projectObjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting project from database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Deleted project %s", id)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /project/stats - ProjectStats returns stats about projects
func ProjectStats(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Aggregate project stats
	stats, err := database.AggregateProjectStats(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating project stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// GET /project/:id - GetProject returns a project by ID
func GetProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the id from the request
	id := ctx.Param("id")

	// Convert project ID string to ObjectID
	projectObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Get the project from the database
	project, err := database.FindProjectById(state.Db, ctx, &projectObjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting project from database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, project)
}

// GET /project/count - GetProjectCount returns the number of projects in the database
func GetProjectCount(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	var count int64

	// Get the options from the database
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options from database: " + err.Error()})
		return
	}

	// If the tracks option is enabled, get the project count for the judge's track
	if options.JudgeTracks && judge.Track != "" {
		// Get the project from the database
		count, err := database.CountTrackProjects(state.Db, judge.Track)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting project count from database: " + err.Error()})
			return
		}

		// Send OK
		ctx.JSON(http.StatusOK, gin.H{"count": count})
		return
	}

	// Get the project from the database
	count, err = database.CountProjectDocuments(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting project count from database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

// PUT /project/hide/:id - Sets the active status of a project
func HideProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get ID from URL
	id := ctx.Param("id")

	// Get request body
	var hideReq util.HideRequest
	err := ctx.BindJSON(&hideReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Convert project ID string to ObjectID
	projectObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Update the project in the database
	err = database.SetProjectActive(state.Db, ctx, &projectObjectId, !hideReq.Hide)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project in database: " + err.Error()})
		return
	}

	// Send OK
	action := "Unhid"
	if hideReq.Hide {
		action = "Hid"
	}
	state.Logger.AdminLogf("%s project %s", action, id)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /project/hide - HideSelectedProjects hides selected projects
func HideSelectedProjects(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get request body
	var hideReq util.HideSelectedRequest
	err := ctx.BindJSON(&hideReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Update the projects in the database
	err = database.SetProjectsActive(state.Db, ctx, hideReq.Items, !hideReq.Hide)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating projects in database: " + err.Error()})
		return
	}

	// Send OK
	action := "Unhid"
	if hideReq.Hide {
		action = "Hid"
	}
	state.Logger.AdminLogf("%s %d projects", action, len(hideReq.Items))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// PUT /project/prioritize/:id - PrioritizeProject prioritizes a project
func PrioritizeProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get ID from URL
	id := ctx.Param("id")

	// Get request body
	var priReq util.PrioritizeRequest
	err := ctx.BindJSON(&priReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Convert project ID string to ObjectID
	projectObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Update the project in the database
	err = database.SetProjectPrioritized(state.Db, &projectObjectId, priReq.Prioritize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project in database: " + err.Error()})
		return
	}

	// Send OK
	action := "Prioritized"
	if !priReq.Prioritize {
		action = "Unprioritized"
	}
	state.Logger.AdminLogf("%s project %s", action, id)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

func PrioritizeSelectedProjects(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request
	var priReq util.PrioritizeSelectedRequest
	err := ctx.BindJSON(&priReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Update the projects in the database
	err = database.SetProjectsPrioritized(state.Db, ctx, priReq.Items, priReq.Prioritize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating projects in database: " + err.Error()})
		return
	}

	// Send OK
	action := "Prioritized"
	if !priReq.Prioritize {
		action = "Unprioritized"
	}
	state.Logger.AdminLogf("%s %d projects", action, len(priReq.Items))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /project/reassign - ReassignProjectNums reassigns project numbers
func ReassignProjectNums(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	err := funcs.ReassignNums(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reassigning project numbers: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Reassigned project numbers")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GetChallenges returns the set of all challenges from the database
func GetChallenges(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the challenges from the database
	challenges, err := database.GetChallenges(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting challenges from database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, challenges)
}

// PUT /project/move/:id - Move a project to a different group
func MoveProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the ID from the URL
	rawId := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(rawId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Get the request object
	var moveReq util.MoveRequest
	err = ctx.BindJSON(&moveReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Run the remaining actions in a transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Get options from database
		options, err := database.GetOptions(state.Db, sc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
			return err
		}

		// Make sure valid group
		if moveReq.Group < 0 || moveReq.Group >= options.NumGroups {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group number"})
			return errors.New("invalid group number")
		}

		// Get max number in group
		currMaxNum, err := database.GetGroupMaxNum(state.Db, sc, moveReq.Group)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number: " + err.Error()})
			return err
		}
		if currMaxNum == -1 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number"})
			return errors.New("error getting group max number")
		}

		// Get maximum number for the given group
		maxGroupNum := int64(0)
		for i, size := range options.GroupSizes {
			maxGroupNum += size
			if int64(i) == moveReq.Group {
				break
			}
		}

		// Make sure the group has space (check for max group number and make sure it doesn't exceed)
		if moveReq.Group < options.NumGroups-1 && currMaxNum >= maxGroupNum {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "group is full, cannot move project to provided group"})
			return err
		}

		// Move the project to the new group
		err = database.SetProjectNum(state.Db, sc, id, currMaxNum+1, moveReq.Group)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error moving project group: " + err.Error()})
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Send OK
	state.Logger.AdminLogf("Moved project %s to group %d", id.Hex(), moveReq.Group)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /project/move - Move selected projects to a different group
func MoveSelectedProjects(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request object
	var moveReq util.MoveSelectedRequest
	err := ctx.BindJSON(&moveReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Run the remaining actions in a transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Get options from database
		options, err := database.GetOptions(state.Db, sc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
			return err
		}

		// Make sure valid group
		if moveReq.Group < 0 || moveReq.Group >= options.NumGroups {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group number"})
			return errors.New("invalid group number")
		}

		// Get max number in group
		currMaxNum, err := database.GetGroupMaxNum(state.Db, ctx, moveReq.Group)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number: " + err.Error()})
			return err
		}
		if currMaxNum == -1 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number"})
			return errors.New("error getting group max number")
		}

		// Get maximum number for the given group
		maxGroupNum := int64(0)
		for i, size := range options.GroupSizes {
			maxGroupNum += size
			if int64(i) == moveReq.Group {
				break
			}
		}

		// Check if the group has space (check for max group number and make sure it doesn't exceed)
		if moveReq.Group < options.NumGroups-1 && currMaxNum+int64(len(moveReq.Items)) > maxGroupNum {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "group is full, cannot move projects to provided group"})
			return errors.New("group is full, cannot move projects to provided group")
		}

		// Move the projects to the new group
		err = database.SetProjectsNum(state.Db, ctx, moveReq.Items, currMaxNum+1, moveReq.Group)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error moving projects group: " + err.Error()})
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Send OK
	state.Logger.AdminLogf("Moved %d projects to group %d", len(moveReq.Items), moveReq.Group)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /project/balance-groups - BalanceProjectGroups balances project groups
func BalanceProjectGroups(ctx *gin.Context) {
	// TODO: Write this; its harder than i imagined oops
	// // Get the database from the context
	// db := ctx.MustGet("db").(*mongo.Database)

	// // Get the options from the database
	// options, err := database.GetOptions(state.Db, ctx)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options from database: " + err.Error()})
	// 	return
	// }

	// // Get all projects
	// projects, err := database.FindAllProjects(state.Db, ctx)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects from database: " + err.Error()})
	// 	return
	// }

	// // Get the project count per group
	// groupCounts := make(map[int64]int, options.NumGroups)
	// for _, proj := range projects {
	// 	groupCounts[proj.Group]++
	// }

	// // Get the average project count per group
	// avg := len(projects) / int(options.NumGroups)
	// rem := len(projects) % int(options.NumGroups)

	// // Figure out which groups need to be adjusted
	// incrGroups := make(map[int64]int, 0)
	// decrGroups := make(map[int64]int, 0)
	// for group, count := range groupCounts {
	// 	comp := avg
	// 	if rem > 0 {
	// 		comp++
	// 		rem--
	// 	}

	// 	if count > comp {
	// 		decrGroups[group] = count - comp
	// 	} else if count < comp {
	// 		incrGroups[group] = comp - count
	// 	}
	// }

	// // Move projects from groups with more projects to groups with fewer projects
	// homeless := make([]*models.Project, 0)
	// for group, count := range decrGroups {
	// 	// Put the last few projects in the group into the homeless list
	// 	// We can loop backwards through the projects bc they are in order of group
	// 	for i := len(projects) - 1; i >= 0; i-- {
	// 		if projects[i].Group == group {
	// 			homeless = append(homeless, projects[i])
	// 			projects = append(projects[:i], projects[i+1:]...)
	// 			count--
	// 			if count <= 0 {
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// // Find homes for the homeless
	// for _, proj := range homeless {

	// }
}

// GET /admin/log - Returns the log file
func GetLog(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the log file
	log := state.Logger.Get()

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"log": log})
}
