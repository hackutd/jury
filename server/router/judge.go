package router

import (
	"errors"
	"net/http"
	"server/database"
	"server/funcs"
	"server/judging"
	"server/models"
	"server/ranking"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GET /judge - Endpoint to get the judge from the token
func GetJudge(ctx *gin.Context) {
	// Get the judge from the context (See middleware.go)
	judge := ctx.MustGet("judge").(*models.Judge)

	// Send Judge
	ctx.JSON(http.StatusOK, judge)
}

// POST /judge/new - Endpoint to add a single judge
func AddJudge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the request
	var judgeReq models.AddJudgeRequest
	err := ctx.BindJSON(&judgeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure all required request fields are defined
	if judgeReq.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if judgeReq.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// Determine group judge should go in
	group := int64(0)
	if judgeReq.Track == "" {
		group, err = database.GetMinJudgeGroup(state.Db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judge group: " + err.Error()})
			return
		}
	}

	// Create the judge
	judge := models.NewJudge(judgeReq.Name, judgeReq.Email, judgeReq.Track, judgeReq.Notes, group)

	// Send email if no_send is false
	if !judgeReq.NoSend {
		// Get hostname from request
		hostname := util.GetFullHostname(ctx)

		// Make sure email is right
		if !funcs.CheckEmail(judge.Email) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}

		// Send email to judge
		err = funcs.SendJudgeEmail(judge, hostname)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error sending judge email: " + err.Error()})
			return
		}
	}

	// Insert the judge into the database
	err = database.InsertJudge(state.Db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OK
	track := judge.Track
	if track == "" {
		track = "general"
	}
	state.Logger.AdminLogf("Added judge %s (%s), track: %s", judge.Name, judge.Email, track)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type LoginJudgeRequest struct {
	Code string `json:"code"`
}

// POST /judge/login - Endpoint to login a judge
func LoginJudge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge code from the request
	var loginReq LoginJudgeRequest
	err := ctx.BindJSON(&loginReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Find judge by code
	judge, err := database.FindJudgeByCode(state.Db, loginReq.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding judge in database: " + err.Error()})
		return
	}
	if judge == nil {
		state.Logger.JudgeLogf(nil, "Invalid judge log in attempt with code %s", loginReq.Code)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	// Generate random 16-character token for judge
	token, err := util.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token: " + err.Error()})
		return
	}

	// Update judge in database with new token
	judge.Token = token
	err = database.UpdateJudge(state.Db, judge, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.JudgeLogf(judge, "Log in, assigned token %s", token)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// POST /judge/auth - Check to make sure a judge is authenticated
func JudgeAuthenticated(ctx *gin.Context) {
	// This route will run the middleware first, and if the middleware
	// passes, then that means the judge is authenticated

	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/csv - Endpoint to add judges from a CSV file
func AddJudgesCsv(ctx *gin.Context) {
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
	judges, err := funcs.ParseJudgeCSV(string(content), hasHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing CSV file: " + err.Error()})
		return
	}

	// Get hostname from request
	hostname := util.GetFullHostname(ctx)

	// Check all judge emails
	for _, judge := range judges {
		if !funcs.CheckEmail(judge.Email) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email: " + judge.Email})
			return
		}
	}

	// Send emails to all judges
	for _, judge := range judges {
		err = funcs.SendJudgeEmail(judge, hostname)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error sending judge " + judge.Name + " email: " + err.Error()})
			return
		}
	}

	// Determine group numbers for the new judges
	groups, err := database.GetNextNJudgeGroups(state.Db, ctx, len(judges), false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judge groups: " + err.Error()})
		return
	}

	// Assign group numbers to the new judges
	for i, judge := range judges {
		judge.Group = groups[i]
	}

	// Insert judges into the database
	err = database.InsertJudges(state.Db, judges)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting judges into database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Added %d judges from CSV", len(judges))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /judge/welcome - Endpoint to check if a judge has read the welcome message
func CheckJudgeReadWelcome(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Send OK
	if judge.ReadWelcome {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}

// POST /judge/welcome - Endpoint to set a judge's readWelcome field to true
func SetJudgeReadWelcome(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Set judge's readWelcome field to true
	judge.ReadWelcome = true

	// Update judge in database
	err := database.UpdateJudge(state.Db, judge, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.JudgeLogf(judge, "Read welcome message")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /judge/list - Endpoint to get a list of all judges
func ListJudges(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get judges from database
	judges, err := database.FindAllJudges(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding judges in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, judges)
}

// GET /judge/stats - Endpoint to get stats about the judges
func JudgeStats(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Aggregate judge stats
	stats, err := database.AggregateJudgeStats(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating judge stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// DELETE /judge/:id - Endpoint to delete a judge
func DeleteJudge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge ID from the URL
	judgeId := ctx.Param("id")

	// Convert judge ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(judgeId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Delete judge from database
	err = database.DeleteJudgeById(state.Db, judgeObjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting judge from database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Deleted judge %s", judgeId)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/next - Endpoint to get the next project for a judge
func GetNextJudgeProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// If the judge already has a next project, return that project
	if judge.Current != nil {
		ctx.JSON(http.StatusOK, gin.H{"project_id": judge.Current.Hex()})
		return
	}

	// Otherwise, get the next project for the judge
	project_int, err := database.WithTransactionItem(state.Db, func(sc mongo.SessionContext) (interface{}, error) {
		// Get options
		options, err := database.GetOptions(state.Db, sc)
		if err != nil {
			return nil, errors.New("error getting options: " + err.Error())
		}

		// If the clock is paused, return an empty object
		// This is to ensure that no projects are gotten if the clock is paused
		state.Clock.Mutex.Lock()
		if !state.Clock.State.Running || options.Deliberation {
			return nil, nil
		}
		state.Clock.Mutex.Unlock()

		return judging.PickNextProject(state.Db, judge, sc, state.Comps)
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error picking next project: " + err.Error()})
		return
	}

	// If there is no next project, return an empty object
	if project_int == nil {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	// Cast the project to a Project object
	project := project_int.(*models.Project)

	// Update judge and project
	err = database.UpdateAfterPicked(state.Db, project, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating next project in database: " + err.Error()})
		return
	}

	// Send OK and project ID
	state.Logger.JudgeLogf(judge, "Picked new project %s (%s)", project.Name, project.Id.Hex())
	ctx.JSON(http.StatusOK, gin.H{"project_id": project.Id.Hex()})
}

// GET /judge/projects - Endpoint to get a list of projects that a judge has seen
func GetJudgeProjects(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Return the judge's seen projects list
	ctx.JSON(http.StatusOK, judge.SeenProjects)
}

type JudgedProjectWithUrl struct {
	models.JudgedProject
	Url string `bson:"url" json:"url"`
}

func addUrlToJudgedProject(project *models.JudgedProject, url string) *JudgedProjectWithUrl {
	return &JudgedProjectWithUrl{
		JudgedProject: models.JudgedProject{
			ProjectId:   project.ProjectId,
			Name:        project.Name,
			Location:    project.Location,
			Description: project.Description,
			Notes:       project.Notes,
			Starred:     project.Starred,
		},
		Url: url,
	}
}

// GET /judge/project/:id - Gets a project that's been judged by ID
func GetJudgedProject(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the project ID from the URL
	projectId := ctx.Param("id")

	// Search through the judge seen projects for the project ID
	for _, p := range judge.SeenProjects {
		if p.ProjectId.Hex() == projectId {
			// Add URL to judged project
			proj, err := database.FindProjectById(state.Db, ctx, &p.ProjectId)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting project url: " + err.Error()})
				return
			}
			jpWithUrl := addUrlToJudgedProject(&p, proj.Url)

			// Parse and send JSON
			ctx.JSON(http.StatusOK, jpWithUrl)
			return
		}
	}

	// Send bad request bc project ID invalid
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
}

type SkipRequest struct {
	Reason string `json:"reason"`
}

// POST /judge/skip - Endpoint to skip a project
func JudgeSkip(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the skip reason from the request
	var skipReq SkipRequest
	err := ctx.BindJSON(&skipReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Skipped project ID
	id := judge.Current.Hex()

	// Check if judging is paused
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}
	state.Clock.Mutex.Lock()
	newProj := state.Clock.State.Running && !options.Deliberation
	state.Clock.Mutex.Unlock()

	// Skip the project
	err = judging.SkipCurrentProject(state.Db, judge, state.Comps, skipReq.Reason, newProj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OK
	if skipReq.Reason == "busy" {
		state.Logger.JudgeLogf(judge, "Skipped busy project %s", id)
	} else {
		state.Logger.JudgeLogf(judge, "Flagged project %s due to %s", id, skipReq.Reason)
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/hide/:id - Endpoint to hide a judge
func HideJudge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get ID from URL
	id := ctx.Param("id")

	// Get ID from body
	var hideReq util.HideRequest
	err := ctx.BindJSON(&hideReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Convert ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Hide judge in database
	err = database.SetJudgeActive(state.Db, &judgeObjectId, !hideReq.Hide)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error hiding judge in database: " + err.Error()})
		return
	}

	// Send OK
	action := "Unhid"
	if hideReq.Hide {
		action = "Hid"
	}
	state.Logger.AdminLogf("%s judge %s", action, id)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/hide - Endpoint to hide selected judges
func HideSelectedJudges(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request object
	var hideReq util.HideSelectedRequest
	err := ctx.BindJSON(&hideReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Hide judges in database
	err = database.SetJudgesActive(state.Db, hideReq.Items, !hideReq.Hide)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error hiding judges in database: " + err.Error()})
		return
	}

	// Send OK
	action := "Unhid"
	if hideReq.Hide {
		action = "Hid"
	}
	state.Logger.AdminLogf("%s %d judges", action, len(hideReq.Items))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// PUT /judge/:id - Endpoint to edit a judge
func EditJudge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the body content
	var judgeReq models.AddJudgeRequest
	err := ctx.BindJSON(&judgeReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the judge ID from the path
	judgeId := ctx.Param("id")

	// Convert ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(judgeId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Edit judge in database
	err = database.UpdateJudgeBasicInfo(state.Db, &judgeObjectId, &judgeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error editing judge in database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Edited judge %s: %s", judgeId, util.StructToString(judgeReq))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type JudgeScoreRequest struct {
	Notes   string `json:"notes"`
	Starred bool   `json:"starred"`
}

// POST /judge/score - Endpoint to finish judging a project
func JudgeScore(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the request object
	var scoreReq JudgeScoreRequest
	err := ctx.BindJSON(&scoreReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the options and return error if deliberations
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}
	if options.Deliberation {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot score due to deliberation mode being enabled"})
		return
	}

	projId := judge.Current.Hex()

	// Wrap the database transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Get the options
		options, err := database.GetOptions(state.Db, sc)
		if err != nil {
			return errors.New("error getting options: " + err.Error())
		}

		// Get the project from the database
		project, err := database.FindProjectById(state.Db, sc, judge.Current)
		if err != nil {
			return errors.New("error finding project in database: " + err.Error())
		}

		// Create the judged project object
		judgedProject := models.JudgeProjectFromProject(project, scoreReq.Notes, scoreReq.Starred)

		// If groups are enabled and auto switch, move the judge to the next group conditionally
		if options.MultiGroup && options.SwitchingMode == "auto" {
			err = judging.MoveJudgeGroup(state.Db, sc, judge, options)
			if err != nil {
				return errors.New("error moving judge group: " + err.Error())
			}
		}

		// Update the judge and project
		err = database.UpdateAfterSeen(state.Db, sc, judge, judgedProject)
		if err != nil {
			return errors.New("error storing scores in database: " + err.Error())
		}

		// Reset list of skipped projects due to busy status
		err = database.ResetBusyProjectListForJudge(state.Db, sc, judge)
		if err != nil {
			return errors.New("error resetting busy project list in database: " + err.Error())
		}

		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OK
	starred := ""
	if scoreReq.Starred {
		starred = " and starred project"
	}
	state.Logger.JudgeLogf(judge, "Finished judging project %s%s", projId, starred)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type RankRequest struct {
	Ranking []primitive.ObjectID `json:"ranking"`
}

// POST /judge/rank - Update the judge's ranking of projects
func JudgeRank(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the request object
	var rankReq RankRequest
	err := ctx.BindJSON(&rankReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the options and return error if deliberations
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}
	if options.Deliberation {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot score due to deliberation mode being enabled"})
		return
	}

	// Old ranks
	oldRanks := util.RankingToString(judge.Rankings)

	// Calculate diff of scores
	diff := ranking.CalculateScoreDiff(rankReq.Ranking, judge.Rankings)

	if len(*diff) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		return
	}

	// Wrap in transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Update the judge's ranking
		err = database.UpdateJudgeRanking(state.Db, judge, rankReq.Ranking)
		if err != nil {
			return errors.New("error updating judge ranking in database: " + err.Error())
		}

		// Update projects based on diff
		err = database.UpdateProjectScores(state.Db, sc, diff)
		if err != nil {
			return errors.New("error updating project scores in database: " + err.Error())
		}

		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OK
	state.Logger.JudgeLogf(judge, "Updated rankings from %s to %s", oldRanks, util.RankingToString(rankReq.Ranking))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type JudgeStarRequest struct {
	Starred bool `json:"starred"`
}

// PUT /judge/star/:id - Update the judge's star status on a judged project
func JudgeStar(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the project id from the URL
	rawProjectId := ctx.Param("id")
	projectId, err := primitive.ObjectIDFromHex(rawProjectId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Get the request object
	var starReq JudgeStarRequest
	err = ctx.BindJSON(&starReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the options and return error if deliberations
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}
	if options.Deliberation {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot score due to deliberation mode being enabled"})
		return
	}

	// If the project isn't in the judge's seen projects, return an error
	index := util.FindSeenProjectIndex(judge, projectId)
	if index == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge hasn't seen project or project is invalid"})
		return
	}

	// Wrap in transaction
	err = database.WithTransaction(state.Db, func(sc mongo.SessionContext) error {
		// Get the options
		options, err := database.GetOptions(state.Db, sc)
		if err != nil {
			return errors.New("error getting options: " + err.Error())
		}

		// Update the judge's object for the project
		err = database.UpdateJudgeStars(state.Db, sc, judge.Id, index, starReq.Starred)
		if err != nil {
			return errors.New("error updating judge starred in database: " + err.Error())
		}

		// If the judge is in a track, update that track's comparisons
		// Otherwise, update the project's starred status
		if options.JudgeTracks && judge.Track != "" {
			err = database.UpdateProjectTrackStars(state.Db, sc, projectId, judge.Track, starReq.Starred)
		} else {
			err = database.UpdateProjectStars(state.Db, sc, projectId, starReq.Starred)
		}
		if err != nil {
			return errors.New("error updating project starred status in database: " + err.Error())
		}

		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OK
	starred := "Removed"
	if starReq.Starred {
		starred = "Added"
	}
	state.Logger.JudgeLogf(judge, "%s star for project %s", starred, rawProjectId)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/break - Allows a judge to take a break and free up their current project
func JudgeBreak(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Error if the judge doesn't have a current project
	if judge.Current == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge doesn't have a current project"})
		return
	}

	// Basically skip the project for the judge
	err := judging.SkipCurrentProject(state.Db, judge, state.Comps, "break", false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error skipping project: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.JudgeLogf(judge, "Took a break from project %s", judge.Current.Hex())
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type UpdateNotesRequest struct {
	Notes string `json:"notes"`
}

// PUT /judge/notes/:id - Update the notes of a judge
func JudgeUpdateNotes(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the project ID from the URL
	rawProjectId := ctx.Param("id")
	projectId, err := primitive.ObjectIDFromHex(rawProjectId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	// Get the request object
	var notesReq UpdateNotesRequest
	err = ctx.BindJSON(&notesReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// If the project isn't in the judge's seen projects, return an error
	index := util.FindSeenProjectIndex(judge, projectId)
	if index == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge hasn't seen project or project is invalid"})
		return
	}

	// Update that specific index of the seen projects array
	judge.SeenProjects[index].Notes = notesReq.Notes

	// Update the judge's object for the project
	err = database.UpdateJudgeSeenProjects(state.Db, ctx, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge score in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/reassign - Reassigns judge numbers to all judges that are not in a track
func ReassignJudgeGroups(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get options from database
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Don't do if not enabled
	if !options.MultiGroup {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "multi-group not enabled"})
		return
	}

	// Reassign judge groups
	err = database.PutJudgesInGroups(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reassigning judge groups: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Reassigned judge groups")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// PUT /judge/move/:id - Move a judge to a different group
func MoveJudge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the ID from the URL
	id := ctx.Param("id")

	// Get the request object
	var moveReq util.MoveRequest
	err := ctx.BindJSON(&moveReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Convert ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Move the judge to the new group
	err = database.SetJudgeGroup(state.Db, ctx, &judgeObjectId, moveReq.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error moving judge group: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Moved judge %s to group %d", id, moveReq.Group)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/move - Move selected judges to a different group
func MoveSelectedJudges(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request object
	var moveReq util.MoveSelectedRequest
	err := ctx.BindJSON(&moveReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Move the judges to the new group
	err = database.SetJudgesGroup(state.Db, ctx, moveReq.Items, moveReq.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error moving judges group: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Moved %d judges to group %d", len(moveReq.Items), moveReq.Group)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
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

	// Get options from database
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Make sure valid group
	if moveReq.Group < 0 || moveReq.Group >= options.NumGroups {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group number"})
		return
	}

	// Get max number in group
	currMaxNum, err := database.GetGroupMaxNum(state.Db, ctx, moveReq.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number: " + err.Error()})
		return
	}
	if currMaxNum == -1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number"})
		return
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
		return
	}

	// Move the project to the new group
	err = database.SetProjectNum(state.Db, ctx, id, currMaxNum+1, moveReq.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error moving project group: " + err.Error()})
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

	// Get options from database
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Make sure valid group
	if moveReq.Group < 0 || moveReq.Group >= options.NumGroups {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group number"})
		return
	}

	// Get max number in group
	currMaxNum, err := database.GetGroupMaxNum(state.Db, ctx, moveReq.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number: " + err.Error()})
		return
	}
	if currMaxNum == -1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting group max number"})
		return
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
		return
	}

	// Move the projects to the new group
	err = database.SetProjectsNum(state.Db, ctx, moveReq.Items, currMaxNum+1, moveReq.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error moving projects group: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Moved %d projects to group %d", len(moveReq.Items), moveReq.Group)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type AddJudgeFromQRRequest struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Track string `json:"track"`
}

// POST /judge/qr - Add a judge from a QR code
func AddJudgeFromQR(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request object
	var qrReq AddJudgeFromQRRequest
	err := ctx.BindJSON(&qrReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the options from the database
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Make sure the code is correct
	if qrReq.Track == "" {
		if qrReq.Code != options.QRCode {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid QR code"})
			return
		}
	} else {
		if qrReq.Code != options.TrackQRCodes[qrReq.Track] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid QR code"})
			return
		}
	}

	// Check if the judge already exists
	judge, err := database.FindJudgeByCode(state.Db, qrReq.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding judge in database: " + err.Error()})
		return
	}
	if judge != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge already exists"})
		return
	}

	// Determine group judge should go in
	group := int64(0)
	if qrReq.Track == "" {
		group, err = database.GetMinJudgeGroup(state.Db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judge group: " + err.Error()})
			return
		}
	}

	// Create the judge
	judge = models.NewJudge(qrReq.Name, qrReq.Email, qrReq.Track, "", group)

	// SEND EMAILS =============================

	// Get hostname from request
	hostname := util.GetFullHostname(ctx)

	// Make sure email is right
	if !funcs.CheckEmail(judge.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	// Send email to judge
	err = funcs.SendJudgeEmail(judge, hostname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error sending judge email: " + err.Error()})
		return
	}

	// Insert the judge into the database
	err = database.InsertJudge(state.Db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting judge into database: " + err.Error()})
		return
	}

	// Send OK
	track := ""
	if qrReq.Track != "" {
		track = ", track " + qrReq.Track
	}
	state.Logger.AdminLogf("Added judge %s (%s) from QR code%s", judge.Name, judge.Email, track)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /judge/deliberation - Get the status of the deliberation
func GetDeliberationStatus(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options from the database
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	if options.Deliberation {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}
