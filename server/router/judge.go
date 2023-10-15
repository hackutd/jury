package router

import (
	"net/http"
	"server/crowdbt"
	"server/database"
	"server/models"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GET /judge - Endpoint to get the judge from the token
func GetJudge(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Send Judge
	ctx.JSON(http.StatusOK, judge)
}

// POST /judge/new - Endpoint to add a single judge
func AddJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

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

	// Create the judge
	judge := models.NewJudge(judgeReq.Name, judgeReq.Email, judgeReq.Notes)

	// Get hostname from request
	hostname := util.GetFullHostname(ctx)

	// Send email to judge
	err = util.SendJudgeEmail(judge, hostname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error sending judge email: " + err.Error()})
		return
	}

	// Insert the judge into the database
	err = database.InsertJudge(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type LoginJudgeRequest struct {
	Code string `json:"code"`
}

// POST /judge/login - Endpoint to login a judge
func LoginJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge code from the request
	var loginReq LoginJudgeRequest
	err := ctx.BindJSON(&loginReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Find judge by code
	judge, err := database.FindJudgeByCode(db, loginReq.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding judge in database: " + err.Error()})
		return
	}
	if judge == nil {
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
	err = database.UpdateJudge(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK
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
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

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
	judges, err := util.ParseJudgeCSV(string(content), hasHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing CSV file: " + err.Error()})
		return
	}

	// Get hostname from request
	hostname := util.GetFullHostname(ctx)

	// Send emails to all judges
	for _, judge := range judges {
		err = util.SendJudgeEmail(judge, hostname)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error sending judge " + judge.Name + " email: " + err.Error()})
			return
		}
	}

	// Insert judges into the database
	err = database.InsertJudges(db, judges)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting judges into database: " + err.Error()})
		return
	}

	// Send OK
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
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Set judge's readWelcome field to true
	judge.ReadWelcome = true

	// Update judge in database
	err := database.UpdateJudge(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /judge/list - Endpoint to get a list of all judges
func ListJudges(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get judges from database
	judges, err := database.FindAllJudges(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding judges in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, judges)
}

// GET /judge/stats - Endpoint to get stats about the judges
func JudgeStats(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Aggregate judge stats
	stats, err := database.AggregateJudgeStats(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating judge stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// DELETE /judge/:id - Endpoint to delete a judge
func DeleteJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge ID from the URL
	judgeId := ctx.Param("id")

	// Convert judge ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(judgeId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Delete judge from database
	err = database.DeleteJudgeById(db, judgeObjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting judge from database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/vote - Endpoint for judge to cast a vote on a project
func JudgeVote(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the vote from the request
	var voteReq models.JudgeVote
	err := ctx.BindJSON(&voteReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get both projects from the database
	prevProject, err := database.FindProjectById(db, judge.Prev)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding previous project in database: " + err.Error()})
		return
	}
	nextProject, err := database.FindProjectById(db, judge.Next)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding next project in database: " + err.Error()})
		return
	}

	// If there is no previous project, then this is the first project for the judge
	if prevProject == nil {
		// Get a new project for the judge
		newProject, err := util.PickNextProject(db, judge)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error picking next project: " + err.Error()})
			return
		}
		if newProject != nil {
			err = database.UpdateProjectSeen(db, newProject, judge)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project seen: " + err.Error()})
				return
			}
			judge.Next = &newProject.Id
		}

		// Update the judge in the DB
		judge.Prev = &nextProject.Id
		err = database.UpdateJudge(db, judge)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
			return
		}

		// Send OK
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		return
	}

	// Set winner/loser based on curr_winner
	var winner, loser *models.Project
	if voteReq.CurrWinner {
		winner, loser = nextProject, prevProject
	} else {
		winner, loser = prevProject, nextProject
	}

	// Run the update function
	nAlpha, nBeta, nMuWinner, nSigmaWinner, nMuLoser, nSigmaLoser := crowdbt.Update(judge.Alpha, judge.Beta, winner.Mu, winner.SigmaSq, loser.Mu, loser.SigmaSq)

	// Update the fields
	judge.Alpha = nAlpha
	judge.Beta = nBeta
	winner.Mu = nMuWinner
	winner.SigmaSq = nSigmaWinner
	loser.Mu = nMuLoser
	loser.SigmaSq = nSigmaLoser

	// Update other fields
	judge.Prev = judge.Next
	judge.Votes += 1
	winner.Votes += 1

	// Get new project for judge
	newProject, err := util.PickNextProject(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error picking next project: " + err.Error()})
		return
	}
	if newProject != nil {
		err = database.UpdateProjectSeen(db, newProject, judge)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project seen: " + err.Error()})
			return
		}
		judge.Next = &newProject.Id
	} else {
		judge.Next = nil
	}

	// Perform database transaction to update judge and both projects
	err = database.UpdateAfterVote(db, judge, winner, loser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge and projects in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /judge/ipo - Endpoint for judge to get IPO (Initial Project Offering) to vote on.
// This endpoint will return a field called "initial".
// If this field is true, then this is the first project for the judge and will return a "project_id" field.
// Otherwise, it will return only the "initial" field as false.
func GetJudgeIPO(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// If the judge has a next or prev project, then this is not the initial project
	if judge.Next != nil || judge.Prev != nil {
		ctx.JSON(http.StatusOK, gin.H{"initial": false})
		return
	}

	// Get the next project for the judge and update the judge/project seen in the database
	// If no project, return initial as false
	project, err := util.PickNextProject(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error picking next project: " + err.Error()})
		return
	}
	if project == nil {
		ctx.JSON(http.StatusOK, gin.H{"initial": false})
		return
	}

	err = database.UpdateProjectSeen(db, project, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project seen: " + err.Error()})
		return
	}

	// Update judge
	judge.Next = &project.Id
	err = database.UpdateJudgeNext(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK and project ID
	ctx.JSON(http.StatusOK, gin.H{"initial": true, "project_id": project.Id})
}

// GET /judge/projects - Endpoint to get a list of projects that a judge has seen
func GetJudgeProjects(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Return the judge's seen projects list
	ctx.JSON(http.StatusOK, judge.SeenProjects)
}

// GET /judge/vote/info - Endpoint to get info about the current projects the judge is voting on
func GetVotingProjectInfo(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the current projects the judge is voting on
	prevProject, err := database.FindProjectById(db, judge.Prev)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding previous project in database: " + err.Error()})
		return
	}
	nextProject, err := database.FindProjectById(db, judge.Next)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding next project in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{
		"prev_name":     prevProject.Name,
		"prev_location": prevProject.Location,
		"curr_name":     nextProject.Name,
		"curr_location": nextProject.Location,
	})
}

type UpdateStarsRequest struct {
	ProjectId string `json:"project_id"`
	Stars     int64  `json:"stars"`
}

// POST /judge/stars - Endpoint to update the stars for a single seen project.
// Body: project_id, stars
func UpdateStars(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the stars from the request
	var starsReq UpdateStarsRequest
	err := ctx.BindJSON(&starsReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Search judge's seen projects for the project ID
	for i, p := range judge.SeenProjects {
		if p.ProjectId.Hex() == starsReq.ProjectId {
			judge.SeenProjects[i].Stars = starsReq.Stars
			break
		}
	}

	// Update the stars for the project
	err = database.UpdateJudgeStars(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge stars in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /judge/project/:id - Gets a project that's been judged by ID
func GetJudgedProject(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the project ID from the URL
	projectId := ctx.Param("id")

	// Search through the judge seen projects for the project ID
	for _, p := range judge.SeenProjects {
		if p.ProjectId.Hex() == projectId {
			ctx.JSON(http.StatusOK, p)
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
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the skip reason from the request
	var skipReq SkipRequest
	err := ctx.BindJSON(&skipReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get skipped project from database
	skippedProject, err := database.FindProjectById(db, judge.Next)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding skipped project in database: " + err.Error()})
		return
	}

	// Create a new skip object
	skip, err := models.NewSkip(skippedProject, judge, skipReq.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating skip object: " + err.Error()})
		return
	}

	// Add skipped project to skipped database
	err = database.InsertSkip(db, skip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting skip into database: " + err.Error()})
		return
	}

	// Get a new project for the judge
	newProject, err := util.PickNextProject(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error picking next project: " + err.Error()})
		return
	}
	if newProject != nil {
		// TODO: We should prob have a separate list for skipped projects so judges can go back
		err = database.UpdateProjectSeen(db, newProject, judge)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating project seen: " + err.Error()})
			return
		}
		judge.Next = &newProject.Id
	} else {
		judge.Next = nil
	}

	// Update the judge in the DB
	err = database.UpdateJudge(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/hide - Endpoint to hide a judge
func HideJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get ID from body
	var idReq models.IdRequest
	err := ctx.BindJSON(&idReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Convert ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(idReq.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Hide judge in database
	err = database.SetJudgeHidden(db, &judgeObjectId, true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error hiding judge in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/unhide - Endpoint to unhide a judge
func UnhideJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get ID from body
	var idReq models.IdRequest
	err := ctx.BindJSON(&idReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}
	id := idReq.Id

	// Convert ID string to ObjectID
	judgeObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge ID"})
		return
	}

	// Unhide judge in database
	err = database.SetJudgeHidden(db, &judgeObjectId, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error unhiding judge in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// PUT /judge/:id - Endpoint to edit a judge
func EditJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

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
	err = database.UpdateJudgeBasicInfo(db, &judgeObjectId, &judgeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error editing judge in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}
