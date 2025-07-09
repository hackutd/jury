package router

import (
	"fmt"
	"net/http"
	"server/config"
	"server/database"
	"server/funcs"
	"server/models"
	"server/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginAdminRequest struct {
	Password string `json:"password"`
}

// POST /admin/login - LoginAdmin authenticates an admin
func LoginAdmin(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

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
		state.Logger.AdminLogf("Log in")
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		return
	}

	// Return status Unauthorized if the password does not match
	state.Logger.AdminLogf("Invalid log in attempt with pass %s", req.Password)
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
	// Get the state from the context
	state := GetState(ctx)

	// Aggregate the stats
	stats, err := database.AggregateStats(state.Db, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// GET /admin/stats/:track - GetAdminStats returns stats about the system
func GetAdminTrackStats(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the track from the URL
	track := ctx.Param("track")

	// Aggregate the stats
	stats, err := database.AggregateStats(state.Db, track)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error aggregating stats: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, stats)
}

// GET /admin/clock - GetClock returns the current clock state
func GetClock(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Unlock the clock
	state.Clock.Mutex.Lock()
	defer state.Clock.Mutex.Unlock()

	// Send clock state
	ctx.JSON(http.StatusOK, gin.H{"running": state.Clock.State.Running, "time": state.Clock.State.GetDuration()})
}

// GET /admin/started - Returns true if the clock is running (NOT paused)
func IsClockRunning(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Unlock the clock
	state.Clock.Mutex.Lock()
	defer state.Clock.Mutex.Unlock()

	// Send OK
	if state.Clock.State.Running {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}

// POST /admin/clock/pause - PauseClock pauses the clock
func PauseClock(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Unlock the clock
	state.Clock.Mutex.Lock()
	defer state.Clock.Mutex.Unlock()

	// Pause the clock
	state.Clock.State.Pause()

	// Backup
	err := database.UpdateClockConditional(state.Db, ctx, &state.Clock.State)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Paused clock")
	ctx.JSON(http.StatusOK, gin.H{"clock": state.Clock.State})
}

// POST /admin/clock/unpause - UnpauseClock unpauses the clock
func UnpauseClock(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Unlock the clock
	state.Clock.Mutex.Lock()
	defer state.Clock.Mutex.Unlock()

	// Unpause the clock
	state.Clock.State.Resume()

	// Backup
	err := database.UpdateClockConditional(state.Db, ctx, &state.Clock.State)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Unpaused clock")
	ctx.JSON(http.StatusOK, gin.H{"clock": state.Clock.State})
}

// POST /admin/clock/reset - ResetClock resets the clock
func ResetClock(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Unlock the clock
	state.Clock.Mutex.Lock()
	defer state.Clock.Mutex.Unlock()

	// Unpause the clock
	state.Clock.State.Resume()

	// Backup
	err := database.UpdateClockConditional(state.Db, ctx, &state.Clock.State)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating clock: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Reset clock")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/clock/backup - BackupClock backups the clock
func BackupClock(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Unlock the clock
	state.Clock.Mutex.Lock()
	defer state.Clock.Mutex.Unlock()

	// Backup the clock
	database.UpdateClock(state.Db, &state.Clock.State)

	// Send OK
	state.Logger.AdminLogf("Backed up clock manually")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/options - sets the options (except for clock and num groups)
func SetOptions(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options from the request
	var options models.OptionalOptions
	err := ctx.BindJSON(&options)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Save the options in the database
	err = database.UpdateOptions(state.Db, ctx, &options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving options: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Updated options: %s", util.StructToStringWithoutNils(options))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

type ResetDatabaseReq struct {
	Type string `json:"type"` // Possible types: all, projects, judges, judging-data, rankings
}

// POST /admin/reset - ResetDatabase resets parts of the database
func ResetDatabase(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get type from request
	var req ResetDatabaseReq
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	fmt.Println(req.Type)

	// Drop based on type
	switch req.Type {
	case "all":
		err = database.DropAll(state.Db)
	case "projects":
		err = database.DropProjects(state.Db)
		if err == nil {
			err = database.DropJudgingData(state.Db)
		}
	case "judges":
		err = database.DropJudges(state.Db)
		if err == nil {
			err = database.DropJudgingData(state.Db)
		}
	case "judging-data":
		database.DropJudgingData(state.Db)
	case "rankings":
		database.DropRankings(state.Db)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid type for reset"})
		return
	}

	// Check for errors
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error resetting database: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Reset database: " + req.Type)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /admin/flags - returns all flags
func GetFlags(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get all the flags
	flags, err := database.FindAllFlags(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting flags: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, flags)
}

// GET /admin/options - returns all options
func GetOptions(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, options)
}

// POST /admin/export/judges - ExportJudges exports all judges to a CSV
func ExportJudges(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get all the judges
	judges, err := database.FindAllJudges(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := funcs.CreateJudgeCSV(judges)

	// Send CSV
	state.Logger.AdminLogf("Exported judges to CSV")
	funcs.AddCsvData("judges", csvData, ctx)
}

// POST /admin/export/projects - ExportProjects exports all projects to a CSV
func ExportProjects(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get all the projects
	projects, err := database.FindAllProjects(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting projects: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := funcs.CreateProjectCSV(projects)

	// Send CSV
	state.Logger.AdminLogf("Exported projects to CSV")
	funcs.AddCsvData("projects", csvData, ctx)
}

// POST /admin/export/challenges - ExportProjectsByChallenge exports all projects to a zip file, with CSVs each
// containing projects that only belong to a single challenge
func ExportProjectsByChallenge(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get all the projects
	projects, err := database.FindAllProjects(state.Db, ctx)
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
	state.Logger.AdminLogf("Exported projects by challenge to CSV")
	funcs.AddZipFile("projects", zipData, ctx)
}

// POST /admin/export/rankings - ExportRankings exports the rankings of each judge as a CSV
func ExportRankings(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get all the judges
	judges, err := database.FindAllJudges(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting judges: " + err.Error()})
		return
	}

	// Create the CSV
	csvData := funcs.CreateJudgeRankingCSV(judges)

	// Send CSV
	state.Logger.AdminLogf("Exported rankings to CSV")
	funcs.AddCsvData("rankings", csvData, ctx)
}

// GET /admin/timer - GetJudgingTimer returns the judging timer
func GetJudgingTimer(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"judging_timer": options.JudgingTimer})
}

// POST /admin/groups/num - SetNumGroups sets the number of groups
func SetNumGroups(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

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
	err = database.UpdateNumGroups(state.Db, *req.NumGroups)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving num groups to options: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Set num groups to %d", req.NumGroups)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/groups/sizes - SetGroupSizes sets the sizes of each group
func SetGroupSizes(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

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
	err = database.UpdateGroupSizes(state.Db, *req.GroupSizes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving group sizes to options: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Set group sizes to %s", util.StructToStringWithoutNils(req))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/groups/swap - SwapJudgeGroups increments the group numbers of all judges
func SwapJudgeGroups(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Swap the groups (increment the group number of each judge)
	err := funcs.IncrementJudgeGroupNum(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error swapping groups: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Swapped judge groups")
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// DELETE /admin/flag/:id - RemoveFlag deletes a flag
func RemoveFlag(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the flag ID from the URL
	id := ctx.Param("id")
	flagId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing flag ID: " + err.Error()})
		return
	}

	// Delete the flag
	err = database.DeleteFlag(state.Db, ctx, &flagId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting flag: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Deleted flag %s", id)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/qr - generates a new QR code
func GenerateQRCode(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Generate QR code
	token, err := util.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating QR code: " + err.Error()})
		return
	}

	// Save the QR code
	err = database.UpdateQRCode(state.Db, ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving QR code: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Generated QR code")
	ctx.JSON(http.StatusOK, gin.H{"qr_code": token})
}

// POST /admin/qr/:track - generates a new QR code for a track
func GenerateTrackQRCode(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the track from the URL
	track := ctx.Param("track")

	// Generate QR code
	token, err := util.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error generating QR code: " + err.Error()})
		return
	}

	// Save the QR code
	err = database.UpdateTrackQRCode(state.Db, ctx, track, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error saving QR code: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Generated QR code for track %s", track)
	ctx.JSON(http.StatusOK, gin.H{"qr_code": token})
}

// GET /admin/qr - GetQRCode returns the QR code
func GetQRCode(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the QR code
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"qr_code": options.QRCode})
}

// GET /admin/qr/:track - GetTrackQRCode returns the QR code for a track
func GetTrackQRCode(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the track from the URL
	track := ctx.Param("track")

	// Get the QR code
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{"qr_code": options.TrackQRCodes[track]})
}

type CheckQRRequest struct {
	Code string `json:"code"`
}

// POST /qr/check - CheckQRCode checks to see if the QR code is right
func CheckQRCode(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request object
	var qrReq CheckQRRequest
	err := ctx.BindJSON(&qrReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the QR code
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK if QR code is right
	if options.QRCode == qrReq.Code {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}

// POST /admin/qr/:track - CheckTrackQRCode checks to see if the track QR code is right
func CheckTrackQRCode(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request object
	var qrReq CheckQRRequest
	err := ctx.BindJSON(&qrReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body: " + err.Error()})
		return
	}

	// Get the track from the URL
	track := ctx.Param("track")

	// Get the QR code
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK if QR code is right
	if options.TrackQRCodes[track] == qrReq.Code {
		ctx.JSON(http.StatusOK, gin.H{"ok": 1})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"ok": 0})
	}
}

type deliberationRequest struct {
	Start bool `json:"start"`
}

// POST /admin/deliberation - SetDeliberation sets the deliberation state
func SetDeliberation(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request
	var req deliberationRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Update the deliberation state
	err = database.UpdateOptions(state.Db, ctx, &models.OptionalOptions{Deliberation: &req.Start})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating deliberation: " + err.Error()})
		return
	}

	// Send OK
	hap := "Started"
	if !req.Start {
		hap = "Stopped"
	}
	state.Logger.AdminLogf("%s deliberation", hap)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// GET /group-info - GetGroupInfo returns the names of the groups and if groups are enabled
func GetGroupInfo(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting options: " + err.Error()})
		return
	}

	// Send OK
	ctx.JSON(http.StatusOK, gin.H{
		"names":   options.GroupNames,
		"enabled": options.MultiGroup,
	})
}

// POST /admin/block-reqs - blocks or unblocks login requests
func SetBlockReqs(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request
	var req models.OptionalOptions
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Make sure the block requests field is valid
	if req.BlockReqs == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid block requests field"})
		return
	}

	// Update the block requests state
	err = database.UpdateOptions(state.Db, ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating block requests: " + err.Error()})
		return
	}

	// Update the limiter
	state.Limiter.Block = *req.BlockReqs

	// Send OK
	state.Logger.AdminLogf("Updated block requests to %t", *req.BlockReqs)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/max-reqs - updates the maximum number of requests per minute
func SetMaxReqs(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the request
	var req models.OptionalOptions
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Make sure the max requests field is valid
	if req.MaxReqPerMin == nil || *req.MaxReqPerMin < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid max requests field"})
		return
	}

	// Update the max requests state
	err = database.UpdateOptions(state.Db, ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating max requests: " + err.Error()})
		return
	}

	// Update the limiter
	state.Limiter.MaxReqPerMin = int(*req.MaxReqPerMin)

	// Send OK
	state.Logger.AdminLogf("Updated max requests to %d", *req.MaxReqPerMin)
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/tracks - updates the tracks to judge
func SetTracks(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not get settings: " + err.Error()})
		return
	}

	// Get the request
	var req models.OptionalOptions
	err = ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Make sure the tracks field is valid
	if req.Tracks == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tracks field"})
		return
	}

	// Make sure track views are equal length. If not, cut or add
	diff := len(*req.Tracks) - len(options.TrackViews)
	changeViews := true
	if diff > 0 { // More tracks = add to views
		for range diff {
			options.TrackViews = append(options.TrackViews, 3) // Default track views = 3
		}
	} else if diff < 0 { // Less tracks = remove from views
		options.TrackViews = options.TrackViews[:len(options.TrackViews)+diff]
	} else {
		changeViews = false
	}

	// Update options
	update := models.OptionalOptions{Tracks: req.Tracks}
	if changeViews {
		update.TrackViews = &options.TrackViews
	}
	err = database.UpdateOptions(state.Db, ctx, &update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating tracks: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Updated tracks to %s", util.StructToStringWithoutNils(req))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}

// POST /admin/track-views - updates the number of views per track
func SetTrackViews(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get the options
	options, err := database.GetOptions(state.Db, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not get settings: " + err.Error()})
		return
	}

	// Get the request
	var req models.OptionalOptions
	err = ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request: " + err.Error()})
		return
	}

	// Make sure the track views field is valid
	if req.TrackViews == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid track views field"})
		return
	}

	// Make sure the track views field is the same length as the tracks
	if len(options.Tracks) != len(*req.TrackViews) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "track views must be the same length as tracks"})
		return
	}

	// Update the options
	err = database.UpdateOptions(state.Db, ctx, &models.OptionalOptions{TrackViews: req.TrackViews})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error updating track views: " + err.Error()})
		return
	}

	// Send OK
	state.Logger.AdminLogf("Updated track views to %s", util.StructToStringWithoutNils(req))
	ctx.JSON(http.StatusOK, gin.H{"ok": 1})
}
