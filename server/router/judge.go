package router

import (
	"net/http"
	"server/database"
	"server/models"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddJudgeRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Notes string `json:"notes"`
}

// POST /judge/new - Endpoint to add a single judge
func AddJudge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the judge from the request
	var judgeReq AddJudgeRequest
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
	ctx.JSON(http.StatusOK, gin.H{"readWelcome": judge.ReadWelcome})
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
	ctx.JSON(http.StatusOK, gin.H{"judges": judges})
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
	ctx.JSON(http.StatusOK, gin.H{"stats": stats})
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

}

// GET /judge/ipo - Endpoint for judge to get IPO (Initial Project Offering) to vote on.
// This endpoint will return a field called "initial".
// If this field is true, then this is the first project for the judge and will return a "project_id" field.
// Otherwise, it will return only the "initial" field as false.
func GetJudgeIPO(ctx *gin.Context) {
	// Get the judge from the context
	judge := ctx.MustGet("judge").(*models.Judge)

	// If the judge has a next or prev project, then this is not the initial project
	if judge.Next != nil || judge.Prev != nil {
		ctx.JSON(http.StatusOK, gin.H{"initial": false})
		return
	}

	// Get the database from the context
	// db := ctx.MustGet("db").(*mongo.Database)

	// TODO: Finish

}
