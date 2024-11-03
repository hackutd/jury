package router

import (
	"net/http"
	"server/config"
	"server/database"
	"server/funcs"
	"server/models"
	"server/ranking"

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

// POST /admin/auth - Checks if an admin is authenticated
func AdminAuthenticated(ctx *gin.Context) {
	// This route will run the middleware first, and if the middleware
	// passes, then that means the admin is authenticated

	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /admin/stats - GetAdminStats returns stats about the system
func GetAdminStats(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Aggregate the stats
	stats, err := database.AggregateStats(db, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// GET /admin/stats/:track - GetAdminStats returns stats about the system
func GetAdminTrackStats(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the track from the URL
	track := ctx.Param("track")

	// Aggregate the stats
	stats, err := database.AggregateStats(db, track)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// GET /admin/clock - GetClock returns the current clock state
func GetClock(ctx *gin.Context) {
	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)

	// Send OK
	sc.Mutex.Lock()

	ctx.JSON(http.StatusOK, gin.H{"running": sc.Clock.Running, "time": sc.Clock.GetDuration()})
	sc.Mutex.Unlock()
}

// POST /admin/clock/pause - PauseClock pauses the clock
func PauseClock(ctx *gin.Context) {
	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	// Pause the clock
	sc.Clock.Pause()

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": sc.Clock})
}

// POST /admin/clock/unpause - UnpauseClock unpauses the clock
func UnpauseClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	// Unpause the clock
	sc.Clock.Resume()

	// Backup
	err := database.UpdateClockConditional(db, ctx, &sc.Clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": sc.Clock})
}

// POST /admin/clock/reset - ResetClock resets the clock
func ResetClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	// Reset the clock
	sc.Clock.Reset()

	// Backup
	err := database.UpdateClockConditional(db, ctx, &sc.Clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": sc.Clock, "ok": 1})
}

// POST /admin/clock/backup - BackupClock backups the clock
func BackupClock(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	// Backup the clock
	database.UpdateClock(db, &sc.Clock)

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"clock": sc.Clock, "ok": 1})
}

type SetClockSyncRequest struct {
	ClockSync bool `json:"clock_sync"`
}

func SetClockSync(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get request
	var req SetClockSyncRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Update the clock sync in options
	err = database.UpdateClockSync(db, ctx, req.ClockSync)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock sync: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /admin/started - Returns true if the clock is running (NOT paused)
func IsClockRunning(ctx *gin.Context) {
	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	// Send OK
	if sc.Clock.Running {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
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

// GET /admin/flags - GetFlags returns all flags
func GetFlags(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the flags
	flags, err := database.FindAllFlags(db)
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
	options, err := database.GetOptions(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, options)
}

// POST /admin/export/judges - ExportJudges exports all judges to a CSV
func ExportJudges(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the judges
	judges, err := database.FindAllJudges(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := funcs.CreateJudgeCSV(judges)

	// Send CSV
	funcs.AddCsvData("judges", csvData, ctx)
}

// POST /admin/export/projects - ExportProjects exports all projects to a CSV
func ExportProjects(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the projects
	projects, err := database.FindAllProjects(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := funcs.CreateProjectCSV(projects)

	// Send CSV
	funcs.AddCsvData("projects", csvData, ctx)
}

// POST /admin/export/challenges - ExportProjectsByChallenge exports all projects to a zip file, with CSVs each
// containing projects that only belong to a single challenge
func ExportProjectsByChallenge(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the projects
	projects, err := database.FindAllProjects(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Create the zip file
	zipData, err := funcs.CreateProjectChallengeZip(projects)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating zip file: " + err.Error()})
		return
	}

	// Send zip file
	funcs.AddZipFile("projects", zipData, ctx)
}

// POST /admin/export/rankings - ExportRankings exports the rankings of each judge as a CSV
func ExportRankings(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the judges
	judges, err := database.FindAllJudges(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := funcs.CreateJudgeRankingCSV(judges)

	// Send CSV
	funcs.AddCsvData("rankings", csvData, ctx)
}

// GET /admin/timer - GetJudgingTimer returns the judging timer
func GetJudgingTimer(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the options
	options, err := database.GetOptions(db, ctx)
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

// POST /admin/timer - SetJudgingTimer sets the judging timer
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
	options, err := database.GetOptions(db, ctx)
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

type SetCategoriesRequest struct {
	Categories []string `json:"categories"`
}

// POST /admin/categories - sets the categories
func SetCategories(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the categories
	var categoriesReq SetCategoriesRequest
	err := ctx.BindJSON(&categoriesReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the categories in the database
	err = database.UpdateCategories(db, categoriesReq.Categories)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving categories: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type MinViewsRequest struct {
	MinViews int `json:"min_views"`
}

// POST /admin/min-views - sets the min views
func SetMinViews(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the views
	var minViewsReq MinViewsRequest
	err := ctx.BindJSON(&minViewsReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
	}

	// Save the min views in the db
	err = database.UpdateMinViews(db, minViewsReq.MinViews)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving min views: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// /GET /admin/score - GetScores returns the calculated scores of all projects
func GetScores(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get all the projects
	projects, err := database.FindAllProjects(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Get all the judges
	judges, err := database.FindAllJudges(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	scores := ranking.CalculateScores(judges, projects)

	// Send OK
	ctx.JSON(http.StatusOK, scores)
}

// /GET /admin/score/<track> - GetTrackScores returns the calculated scores of all projects in a track
func GetTrackScores(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the track from the URL
	track := ctx.Param("track")

	// Get all the projects
	projects, err := database.FindProjectsByTrack(db, ctx, track)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Get all the judges
	judges, err := database.FindJudgesByTrack(db, ctx, track)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	scores := ranking.CalculateScores(judges, projects)

	// Send OK
	ctx.JSON(http.StatusOK, scores)

}

type ToggleTracksRequest struct {
	JudgeTracks bool `json:"judge_tracks"`
}

// POST /admin/tracks/toggle - Toggles the track setting
func ToggleTracks(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the request
	var req ToggleTracksRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the options in the database
	err = database.UpdateJudgeTracks(db, ctx, req.JudgeTracks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type SetTracksRequest struct {
	Tracks []string `json:"tracks"`
}

// POST /admin/tracks - SetTracks sets the tracks
func SetTracks(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the tracks
	var tracksReq SetTracksRequest
	err := ctx.BindJSON(&tracksReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the tracks in the database
	err = database.UpdateTracks(db, ctx, tracksReq.Tracks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving tracks: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type ToggleGroupsRequest struct {
	MultiGroup bool `json:"multi_group"`
}

// POST /admin/groups/toggle - ToggleGroups toggles the multi-group setting
func ToggleGroups(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the request
	var req ToggleGroupsRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Reassign table numbers based on groups
	if req.MultiGroup {
		err = funcs.ReassignNumsByGroup(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reassigning table numbers: " + err.Error()})
			return
		}
	} else {
		err = funcs.ReassignNumsInOrder(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error reassigning table numbers (no groups): " + err.Error()})
			return
		}
	}

	// Save the options in the database
	err = database.UpdateMultiGroup(db, ctx, req.MultiGroup)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type SetNumGroupsRequest struct {
	NumGroups int64 `json:"num_groups"`
}

// POST /admin/groups/num - SetNumGroups sets the number of groups
func SetNumGroups(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// TODO: Wrap in transaction, need to reset groups of everything if changing this

	// Get the request
	var req SetNumGroupsRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the options in the database
	err = database.UpdateNumGroups(db, ctx, req.NumGroups)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type SetGroupSizesRequest struct {
	GroupSizes []int64 `json:"group_sizes"`
}

// POST /admin/groups/sizes - SetGroupSizes sets the group sizes
func SetGroupSizes(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the request
	var req SetGroupSizesRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the options in the database
	err = database.UpdateGroupSizes(db, ctx, req.GroupSizes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/groups/options - SetGroupOptions sets the group options based on the request
func SetGroupOptions(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the request
	var req models.OptionalGroupOptions
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the options in the database
	err = database.UpdateGroupOptions(db, ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/groups/swap - SwapJudgeGroups increments the group numbers of all judges
func SwapJudgeGroups(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Swap the groups (increment the group number of each judge)
	err := funcs.IncrementJudgeGroupNum(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error swapping groups: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}
