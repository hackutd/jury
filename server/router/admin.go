package router

import (
	"fmt"
	"net/http"
	"server/config"
	"server/database"
	"server/models"
	"server/util"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginAdminRequest struct {
	Password string `json:"password"`
}

// POST /admin/login - LoginAdmin authenticates an admin
func LoginAdmin(ctx *gin.Context) {
	// Get the password from the environmental variable
	password := config.GetEnv("JURY_ADMIN_PASSWORD")

	// Get password guess from request
	var req LoginAdminRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Return status OK if the password matches
	if req.Password == password {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		return
	}

	// Return status Unauthorized if the password does not match
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing password field"})
}

// GET /admin/stats - GetAdminStats returns stats about the system
func GetAdminStats(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Aggregate the stats
	stats, err := database.AggregateStats(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)

}

// GET /admin/clock - GetClock returns the current clock state
func GetClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	clock := ctx.MustGet("clock").(*models.ClockState)

	// Save the options in the database
	err := database.UpdateClock(db, clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"running": clock.Running, "time": clock.GetDuration()})
}

// POST /admin/clock/pause - PauseClock pauses the clock
func PauseClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	clock := ctx.MustGet("clock").(*models.ClockState)

	// Pause the clock
	clock.Pause()

	// Save the clock in the database
	err := database.UpdateClock(db, clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving clock: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": clock})
}

// POST /admin/clock/unpause - UnpauseClock unpauses the clock
func UnpauseClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	clock := ctx.MustGet("clock").(*models.ClockState)

	// Unpause the clock
	clock.Resume()

	// Save the clock in the database
	err := database.UpdateClock(db, clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving clock: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": clock})
}

// POST /admin/clock/reset - ResetClock resets the clock
func ResetClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	clock := ctx.MustGet("clock").(*models.ClockState)

	// Reset the clock
	clock.Reset()

	// Save the clock in the database
	err := database.UpdateClock(db, clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving clock: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": clock})
}

// POST /admin/auth - Checks if an admin is authenticated
func AdminAuthenticated(ctx *gin.Context) {
	// This route will run the middleware first, and if the middleware
	// passes, then that means the admin is authenticated

	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/reset - ResetDatabase resets the database
func ResetDatabase(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Reset the database
	err := database.DropAll(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error resetting database: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

func IsClockPaused(ctx *gin.Context) {
	// Get the clock from the context
	clock := ctx.MustGet("clock").(*models.ClockState)

	// Send OK
	if clock.Running {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}

// POST /admin/flags - GetFlags returns all flags
func GetFlags(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the flags
	flags, err := database.FindAllSkips(db, true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting flags: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, flags)
}

func GetOptions(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the options
	options, err := database.GetOptions(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, options)
}

type SetGroupsRequest struct {
	UseGroups bool   `json:"use_groups"`
	RawGroups string `json:"raw_groups"`
}

// POST /admin/options - SetOptions sets the options object
func SetGroups(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get request
	var req SetGroupsRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Parse the groups string; it should be in the form <int> <int>, separated by newlines
	groupsSplit := strings.Split(req.RawGroups, "\n")
	groups := make([]models.Group, 0)

	for i, group := range groupsSplit {
		// Ignore empty lines
		if strings.TrimSpace(group) == "" {
			continue
		}

		groupSplitSplit := strings.Split(group, " ")
		start, err := strconv.Atoi(groupSplitSplit[0])
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error on line %d with start: %s", i+1, err.Error())})
			return
		}

		end, err := strconv.Atoi(groupSplitSplit[1])
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error on line %d with end: %s", i+1, err.Error())})
			return
		}

		groups = append(groups, *models.NewGroup(start, end))
	}

	// Sort groups by start time
	sort.Sort(models.ByStartTime(groups))

	// Check for overlapping groups
	for i := 0; i < len(groups)-1; i++ {
		if groups[i].End >= groups[i+1].Start {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error: group (%d, %d) overlaps with previous group(s)", groups[i+1].Start, groups[i+1].End)})
			return
		}
	}

	// Get the options
	options, err := database.GetOptions(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	options.Groups = groups
	options.UseGroups = req.UseGroups

	// Save the options in the database
	err = database.UpdateOptions(db, options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/export/judges - ExportJudges exports all judges to a CSV
func ExportJudges(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the judges
	judges, err := database.FindAllJudges(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := util.CreateJudgeCSV(judges)

	// Send CSV
	util.AddCsvData("judges", csvData, ctx)
}

// POST /admin/export/projects - ExportProjects exports all projects to a CSV
func ExportProjects(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the projects
	projects, err := database.FindAllProjects(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := util.CreateProjectCSV(projects)

	// Send CSV
	util.AddCsvData("projects", csvData, ctx)
}

// POST /admin/export/challenges - ExportProjectsByChallenge exports all projects to a zip file, with CSVs each
// containing projects that only belong to a single challenge
func ExportProjectsByChallenge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the projects
	projects, err := database.FindAllProjects(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Create the zip file
	zipData, err := util.CreateProjectChallengeZip(projects)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating zip file: " + err.Error()})
		return
	}

	// Send zip file
	util.AddZipFile("projects", zipData, ctx)
}

func GetJudgingTimer(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the options
	options, err := database.GetOptions(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"judging_timer": options.JudgingTimer})
}

type SetJudgingTimerRequest struct {
	JudgingTimer int64 `json:"judging_timer"`
}

func SetJudgingTimer(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get request
	var req SetJudgingTimerRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Get the options
	options, err := database.GetOptions(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	options.JudgingTimer = req.JudgingTimer

	// Save the options in the database
	err = database.UpdateOptions(db, options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}
