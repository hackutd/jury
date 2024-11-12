package router

import (
	"net/http"
	"server/config"
	"server/database"
	"server/funcs"
	"server/logging"
	"server/models"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginAdminRequest struct {
	Password string `json:"password"`
}

// POST /admin/login - LoginAdmin authenticates an admin
func LoginAdmin(ctx *gin.Context) {
	// Get the password from the environmental variable
	password := config.GetEnv("JURY_ADMIN_PASSWORD")

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get password guess from request
	var req LoginAdminRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Return status OK if the password matches
	if req.Password == password {
		logger.AdminLogf("Log in")
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		return
	}

	// Return status Unauthorized if the password does not match
	logger.AdminLogf("Invalid log in attempt with pass %s", req.Password)
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

// POST /admin/clock/pause - PauseClock pauses the clock
func PauseClock(ctx *gin.Context) {
	// Get the clock from the context
	sc := ctx.MustGet("clock").(*models.SafeClock)
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Pause the clock
	sc.Clock.Pause()

	// Send OK
	logger.AdminLogf("Paused clock")
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

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Unpause the clock
	sc.Clock.Resume()

	// Backup
	err := database.UpdateClockConditional(db, ctx, &sc.Clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Unpaused clock")
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

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Reset the clock
	sc.Clock.Reset()

	// Backup
	err := database.UpdateClockConditional(db, ctx, &sc.Clock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Reset clock")
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

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Backup the clock
	database.UpdateClock(db, &sc.Clock)

	// Send OK
	logger.AdminLogf("Backed up clock")
	ctx.JSON(http.StatusOK, gin.H{"clock": sc.Clock, "ok": 1})
}

// POST /admin/options - sets the options (except for clock and num groups)
func SetOptions(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get the options
	var options models.OptionalOptions
	err := ctx.BindJSON(&options)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the options in the database
	err = database.UpdateOptions(db, ctx, &options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Updated options: %s", util.StructToStringWithoutNils(options))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/reset - ResetDatabase resets the database
func ResetDatabase(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Reset the database
	err := database.DropAll(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error resetting database: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("RESET ENTIRE DATABASE")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /admin/flags - returns all flags
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

// GET /admin/options - returns all options
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

// POST /admin/groups/num - SetNumGroups sets the number of groups
func SetNumGroups(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get the request
	var req models.OptionalOptions
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Make sure the number of groups is valid
	if req.NumGroups == nil || *req.NumGroups < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid number of groups"})
		return
	}

	// Update number of groups in the database
	err = database.WithTransaction(db, func(sc mongo.SessionContext) error {
		return database.UpdateNumGroups(db, sc, *req.NumGroups)
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving num groups to options: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Set num groups to %d", req.NumGroups)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/groups/sizes - SetGroupSizes sets the sizes of each group
func SetGroupSizes(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get the request
	var req models.OptionalOptions
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Make sure the group sizes are valid
	if req.GroupSizes == nil || len(*req.GroupSizes) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid group sizes"})
		return
	}

	// Update group sizes in the database
	err = database.WithTransaction(db, func(sc mongo.SessionContext) error {
		return database.UpdateGroupSizes(db, sc, *req.GroupSizes)
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving group sizes to options: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Set group sizes to %s", util.StructToStringWithoutNils(req))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/groups/swap - SwapJudgeGroups increments the group numbers of all judges
func SwapJudgeGroups(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Swap the groups (increment the group number of each judge)
	err := funcs.IncrementJudgeGroupNum(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error swapping groups: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Swapped judge groups")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// DELETE /admin/flag/:id - RemoveFlag deletes a flag
func RemoveFlag(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get the flag ID from the URL
	id := ctx.Param("id")
	flagId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing flag ID: " + err.Error()})
		return
	}

	// Delete the flag
	err = database.DeleteFlag(db, ctx, &flagId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting flag: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Deleted flag %s", id)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/qr - generates a new QR code
func GenerateQRCode(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Generate QR code
	token, err := util.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating QR code: " + err.Error()})
		return
	}

	// Save the QR code
	err = database.UpdateQRCode(db, ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving QR code: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Generated QR code")
	ctx.JSON(http.StatusOK, gin.H{"qr_code": token})
}

// POST /admin/qr/:track - generates a new QR code for a track
func GenerateTrackQRCode(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get the track from the URL
	track := ctx.Param("track")

	// Generate QR code
	token, err := util.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating QR code: " + err.Error()})
		return
	}

	// Save the QR code
	err = database.UpdateTrackQRCode(db, ctx, track, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving QR code: " + err.Error()})
		return
	}

	// Send OK
	logger.AdminLogf("Generated QR code for track %s", track)
	ctx.JSON(http.StatusOK, gin.H{"qr_code": token})
}

// GET /admin/qr - GetQRCode returns the QR code
func GetQRCode(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the QR code
	options, err := database.GetOptions(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"qr_code": options.QRCode})
}

// GET /admin/qr/:track - GetTrackQRCode returns the QR code for a track
func GetTrackQRCode(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the track from the URL
	track := ctx.Param("track")

	// Get the QR code
	options, err := database.GetOptions(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"qr_code": options.TrackQRCodes[track]})
}

type deliberationRequest struct {
	Start bool `json:"start"`
}

// POST /admin/deliberation - SetDeliberation sets the deliberation state
func SetDeliberation(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the logger from the context
	logger := ctx.MustGet("logger").(*logging.Logger)

	// Get the request
	var req deliberationRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Update the deliberation state
	err = database.UpdateDeliberation(db, ctx, req.Start)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating deliberation: " + err.Error()})
		return
	}

	// Send OK
	hap := "Started"
	if !req.Start {
		hap = "Stopped"
	}
	logger.AdminLogf("%s deliberation", hap)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

func GetGroupNames(ctx *gin.Context) {
	// Get the database from the context
	db := ctx.MustGet("db").(*mongo.Database)

	// Get the options
	options, err := database.GetOptions(db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, options.GroupNames)
}
