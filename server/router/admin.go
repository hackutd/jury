package router

import (
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

	// Reset the clock
	state.Clock.State.Reset()

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

// GET /admin/judges - GetJudges returns all judges, including their codes
func GetJudges(ctx *gin.Context) {
	// Get the state from the context
	state := GetState(ctx)

	// Get all judges from the database
	judges, err := database.FindAllJudges(state.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching judges: " + err.Error()})
		return
	}

	// The plan requires ensuring the judge's code is returned.
	// Assuming `database.FindAllJudges` returns a slice of `models.Judge`
	// and the `models.Judge` struct has a `Code` field with a `json:"code"` tag,
	// directly returning the slice will serialize all necessary fields, including the code.
	ctx.JSON(http.StatusOK, judges)
}