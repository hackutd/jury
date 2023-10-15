package router

import (
	"net/http"
	"server/config"
	"server/database"
	"server/models"

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
