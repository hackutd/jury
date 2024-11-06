package router

import (
	"context"
	"log"
	"server/database"
	"server/judging"
	"server/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Creates a new Gin router with all routes defined
func NewRouter(db *mongo.Database) *gin.Engine {
	// Create the router
	router := gin.Default()

	// Get the clock state from the database
	clock := getClockFromDb(db)

	// Create the comparisons object
	comps, err := judging.LoadComparisons(db)
	if err != nil {
		log.Fatalf("error loading projects from the database: %s\n", err.Error())
	}

	// Create track comparisons array
	trackComps, err := judging.LoadTrackComparisons(db)
	if err != nil {
		log.Fatalf("error loading track comparisons from the database: %s\n", err.Error())
	}

	// Add shared variables to router
	router.Use(useVar("db", db))
	router.Use(useVar("clock", clock))
	router.Use(useVar("comps", comps))
	router.Use(useVar("trackComps", trackComps))

	// CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Create router groups for judge and admins
	// This grouping allows us to add middleware to all routes in the group
	// Admins are authenticated with a password and judges are authenticated with a token
	// The default group is for routes that do not require authentication
	judgeRouter := router.Group("/api", AuthenticateJudge())
	adminRouter := router.Group("/api", AuthenticateAdmin())
	defaultRouter := router.Group("/api")

	// ###############################
	// ##### DEFINE ROUTES BELOW #####
	// ###############################

	defaultRouter.GET("/", Heartbeat)

	// Login routes
	defaultRouter.POST("/judge/login", LoginJudge)
	defaultRouter.POST("/admin/login", LoginAdmin)
	judgeRouter.POST("/judge/auth", JudgeAuthenticated)
	adminRouter.POST("/admin/auth", AdminAuthenticated)

	// Admin panel - judges
	adminRouter.POST("/judge/new", AddJudge)
	adminRouter.POST("/judge/csv", AddJudgesCsv)
	adminRouter.GET("/judge/list", ListJudges)
	adminRouter.DELETE("/judge/:id", DeleteJudge)
	adminRouter.PUT("/judge/:id", EditJudge)
	adminRouter.POST("/admin/groups/swap", SwapJudgeGroups)

	// Admin panel - projects
	adminRouter.POST("/project/new", AddProject)
	adminRouter.POST("/project/devpost", AddDevpostCsv)
	adminRouter.POST("/project/csv", AddProjectsCsv)
	adminRouter.GET("/project/list", ListProjects)
	adminRouter.DELETE("/project/:id", DeleteProject)
	// TODO: Add edit project route

	// Admin panel - stats/data
	adminRouter.GET("/admin/stats", GetAdminStats)
	adminRouter.GET("/admin/stats/:track", GetAdminTrackStats)
	adminRouter.GET("/project/stats", ProjectStats)
	adminRouter.GET("/judge/stats", JudgeStats)
	adminRouter.GET("/admin/score", GetScores)
	adminRouter.GET("/admin/score/:track", GetTrackScores)
	adminRouter.GET("/admin/stars", GetStars)
	adminRouter.GET("/admin/stars/:track", GetTrackStars)
	adminRouter.GET("/admin/flags", GetFlags)

	// Admin panel - clock
	adminRouter.GET("/admin/clock", GetClock)
	adminRouter.POST("/admin/clock/pause", PauseClock)
	adminRouter.POST("/admin/clock/unpause", UnpauseClock)
	adminRouter.POST("/admin/clock/sync", SetClockSync)
	adminRouter.POST("/admin/clock/backup", BackupClock)
	defaultRouter.GET("/admin/started", IsClockRunning)

	// Admin panel - options/settings
	adminRouter.POST("/admin/clock/reset", ResetClock)
	adminRouter.POST("/admin/reset", ResetDatabase)
	adminRouter.POST("/project/reassign", ReassignProjectNums)
	adminRouter.POST("/judge/reassign", ReassignJudgeGroups)
	judgeRouter.GET("/admin/timer", GetJudgingTimer)
	adminRouter.POST("/admin/timer", SetJudgingTimer)
	adminRouter.POST("/admin/categories", SetCategories)
	adminRouter.POST("/admin/min-views", SetMinViews)
	adminRouter.GET("/admin/options", GetOptions)
	adminRouter.POST("/admin/tracks/toggle", ToggleTracks)
	adminRouter.POST("/admin/tracks", SetTracks)
	adminRouter.POST("/admin/groups/toggle", ToggleGroups)
	adminRouter.POST("/admin/groups/num", SetNumGroups)
	adminRouter.POST("/admin/groups/sizes", SetGroupSizes)
	adminRouter.POST("/admin/groups/options", SetGroupOptions)

	// Admin panel - exports
	adminRouter.GET("/admin/export/judges", ExportJudges)
	adminRouter.GET("/admin/export/projects", ExportProjects)
	adminRouter.GET("/admin/export/challenges", ExportProjectsByChallenge)
	adminRouter.GET("/admin/export/rankings", ExportRankings)

	// Admin panel - table actions
	adminRouter.POST("/judge/hide", HideJudge)
	adminRouter.POST("/judge/unhide", UnhideJudge)
	adminRouter.POST("/project/hide", HideProject)
	adminRouter.POST("/project/unhide", UnhideProject)
	adminRouter.POST("/project/prioritize", PrioritizeProject)
	adminRouter.POST("/project/unprioritize", UnprioritizeProject)

	// Judging
	judgeRouter.GET("/judge", GetJudge)
	judgeRouter.GET("/judge/welcome", CheckJudgeReadWelcome)
	judgeRouter.POST("/judge/welcome", SetJudgeReadWelcome)
	judgeRouter.GET("/judge/projects", GetJudgeProjects)
	judgeRouter.POST("/judge/next", GetNextJudgeProject)
	judgeRouter.POST("/judge/skip", JudgeSkip)
	judgeRouter.POST("/judge/score", JudgeScore)
	judgeRouter.POST("/judge/rank", JudgeRank)
	judgeRouter.PUT("/judge/star/:id", JudgeStar)
	judgeRouter.PUT("/judge/score", JudgeUpdateScore)
	judgeRouter.POST("/judge/break", JudgeBreak)
	judgeRouter.POST("/judge/notes", JudgeUpdateNotes)
	judgeRouter.GET("/project/:id", GetProject)
	judgeRouter.GET("/project/count", GetProjectCount)
	judgeRouter.GET("/judge/project/:id", GetJudgedProject)
	judgeRouter.GET("/categories", GetCategories)

	// Project expo routes
	defaultRouter.GET("/project/list/public", ListPublicProjects)
	defaultRouter.GET("/challenges", GetChallenges)

	// ######################
	// ##### END ROUTES #####
	// ######################

	// Serve frontend static files
	router.Use(static.Serve("/assets", static.LocalFile("./public/assets", true)))
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.LoadHTMLFiles("./public/index.html")

	// Add no route handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})

	return router
}

// useVar is a middleware that adds a variable to the context
func useVar(key string, v any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(key, v)
		ctx.Next()
	}
}

// getClockFromDb gets the clock state from the database
// and on init will pause the clock
func getClockFromDb(db *mongo.Database) *models.SafeClock {
	// Get the clock state from the database
	options, err := database.GetOptions(db, context.Background())
	if err != nil {
		log.Fatalln("error getting options: " + err.Error())
	}
	clock := options.Clock

	// If the sync clock option is not enabled, return 0 clock
	if !options.ClockSync {
		return models.NewSafeClock(models.NewClockState())
	}

	// Wrap the clock in a mutex
	mut := models.NewSafeClock(&clock)

	return mut
}

// Heartbeat is a simple endpoint to check if the server is running
func Heartbeat(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"ok": 1})
}
