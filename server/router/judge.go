package router

import (
	"fmt"
	"net/http"
	"server/database"
	"server/funcs"
	"server/judging"
	"server/models"
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

	fmt.Println(judgeReq.NoSend)

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

// POST /judge/next - Endpoint to get the next project for a judge
func GetNextJudgeProject(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the comparisons from the context
	comps := ctx.MustGet("comps").(*judging.Comparisons)

	// If the judge already has a next project, return that project
	if judge.Current != nil {
		ctx.JSON(http.StatusOK, gin.H{"project_id": judge.Current.Hex()})
		return
	}

	// Otherwise, get the next project for the judge
	// TODO: This wrapping is a little ridiculous...
	var project *models.Project
	err := database.WithTransaction(db, func(ctx mongo.SessionContext) (interface{}, error) {
		var err error
		project, err = judging.PickNextProject(db, judge, ctx, comps)
		return nil, err
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error picking next project: " + err.Error()})
		return
	}

	// If there is no next project, return an empty object
	if project == nil {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	// Update judge and project
	err = database.UpdateAfterPicked(db, project, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating next project in database: " + err.Error()})
		return
	}

	// Send OK and project ID
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
			Categories:  project.Categories,
			Name:        project.Name,
			Location:    project.Location,
			Description: project.Description,
			Notes:       project.Notes,
		},
		Url: url,
	}
}

// GET /judge/project/:id - Gets a project that's been judged by ID
func GetJudgedProject(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the project ID from the URL
	projectId := ctx.Param("id")

	// Search through the judge seen projects for the project ID
	for _, p := range judge.SeenProjects {
		if p.ProjectId.Hex() == projectId {
			// Add URL to judged project
			proj, err := database.FindProjectById(db, &p.ProjectId)
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
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the comparisons object
	comps := ctx.MustGet("comps").(*judging.Comparisons)

	// Get the skip reason from the request
	var skipReq SkipRequest
	err := ctx.BindJSON(&skipReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Skip the project
	err = judging.SkipCurrentProject(db, judge, comps, skipReq.Reason, true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

type JudgeScoreRequest struct {
	Categories map[string]int `json:"categories"`
}

// POST /judge/score - Endpoint to finish judging a project (give it a score in categories)
func JudgeScore(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the request object
	var scoreReq JudgeScoreRequest
	err := ctx.BindJSON(&scoreReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the project from the database
	project, err := database.FindProjectById(db, judge.Current)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error finding project in database: " + err.Error()})
		return
	}

	// Create the judged project object
	judgedProject := models.JudgeProjectFromProject(project, scoreReq.Categories)

	// Update the project with the score
	err = database.UpdateAfterSeen(db, judge, judgedProject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error storing scores in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type RankRequest struct {
	Ranking []primitive.ObjectID `json:"ranking"`
}

// POST /judge/rank - Update the judge's ranking of projects
func JudgeRank(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the request object
	var rankReq RankRequest
	err := ctx.BindJSON(&rankReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Update the judge's ranking
	err = database.UpdateJudgeRanking(db, judge, rankReq.Ranking)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge ranking in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /judge/break - Allows a judge to take a break and free up their current project
func JudgeBreak(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the comparisons from the context
	comps := ctx.MustGet("comps").(*judging.Comparisons)

	// Error if the judge doesn't have a current project
	if judge.Current == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge doesn't have a current project"})
		return
	}

	// Basically skip the project for the judge
	err := judging.SkipCurrentProject(db, judge, comps, "break", false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error skipping project: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /categories - Endpoint to get the categories
func GetCategories(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get categories from database
	categories, err := database.GetCategories(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting categories: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, categories)
}

type UpdateScoreRequest struct {
	Categories map[string]int     `json:"categories"`
	Project    primitive.ObjectID `json:"project"`
	Initial    bool               `json:"initial"`
}

// PUT /judge/score - Endpoint to update a judge's score for a certain project
func JudgeUpdateScore(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the comparisons object from the context
	comps := ctx.MustGet("comps").(*judging.Comparisons)

	// Get the request object
	var scoreReq UpdateScoreRequest
	err := ctx.BindJSON(&scoreReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Find the index of the project in the judge's seen projects
	index := -1
	for i, p := range judge.SeenProjects {
		if p.ProjectId == scoreReq.Project {
			index = i
			break
		}
	}

	// If the project isn't in the judge's seen projects, return an error
	if index == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge hasn't seen project or project is invalid"})
		return
	}

	// Update that specific index of the seen projects array
	judge.SeenProjects[index].Categories = scoreReq.Categories

	// Update the judge's score for the project
	err = database.UpdateJudgeSeenProjects(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge score in database: " + err.Error()})
		return
	}

	// If this is the initial scoring, update the comparisons array
	comps.UpdateProjectComparisonCount(judge.SeenProjects, scoreReq.Project)

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type UpdateNotesRequest struct {
	Notes   string             `json:"notes"`
	Project primitive.ObjectID `json:"project"`
}

// POST /judge/notes - Update the notes of a judge
func JudgeUpdateNotes(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// Get the request object
	var scoreReq UpdateNotesRequest
	err := ctx.BindJSON(&scoreReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Find the index of the project in the judge's seen projects
	// TODO: Extract to diff function to get rid of repeated code from JudgeUpdateScore
	index := -1
	for i, p := range judge.SeenProjects {
		if p.ProjectId == scoreReq.Project {
			index = i
			break
		}
	}

	// If the project isn't in the judge's seen projects, return an error
	if index == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "judge hasn't seen project or project is invalid"})
		return
	}

	// Update that specific index of the seen projects array
	judge.SeenProjects[index].Notes = scoreReq.Notes

	// Update the judge's object for the project
	err = database.UpdateJudgeSeenProjects(db, judge)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating judge score in database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}
