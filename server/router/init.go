package router

import (
	"context"
	"log"
	"server/database"
	"server/judging"
	"server/logging"
	"server/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Creates a new Gin router with all routes defined
func NewRouter(db *mongo.Database, logger *logging.Logger) *gin.Engine {
	// Create the router
	router := gin.Default()

	// Get the clock state from the database
	clock := getClockFromDb(db)

	// Create the comparisons object
	comps, err := judging.LoadComparisons(db)
	if err != nil {
		log.Fatalf("error loading projects from the database: %s\n", err.Error())
	}

	// Get the limiter from the database
	limiter := getLimiterFromDb(db)

	// Add shared variables to router
	router.Use(useVar("db", db))
	router.Use(useVar("clock", clock))
	router.Use(useVar("comps", comps))
	router.Use(useVar("logger", logger))
	router.Use(useVar("limiter", limiter))

	// Rate limit login requests
	router.Use(rateLimit(limiter))

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
	adminRouter.POST("/admin/qr", GenerateQRCode)
	adminRouter.POST("/admin/qr/:track", GenerateTrackQRCode)
	defaultRouter.GET("/qr", GetQRCode)
	defaultRouter.GET("/qr/:track", GetTrackQRCode)
	defaultRouter.POST("/qr/add", AddJudgeFromQR)

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
	adminRouter.GET("/admin/flags", GetFlags)

	// Admin panel - clock
	adminRouter.GET("/admin/clock", GetClock)
	adminRouter.POST("/admin/clock/pause", PauseClock)
	adminRouter.POST("/admin/clock/unpause", UnpauseClock)
	adminRouter.POST("/admin/clock/backup", BackupClock)
	defaultRouter.GET("/admin/started", IsClockRunning)

	// Admin panel - options/settings
	adminRouter.POST("/admin/clock/reset", ResetClock)
	adminRouter.POST("/admin/reset", ResetDatabase)
	adminRouter.POST("/project/reassign", ReassignProjectNums)
	adminRouter.POST("/project/balance-groups", BalanceProjectGroups)
	adminRouter.POST("/judge/reassign", ReassignJudgeGroups)
	judgeRouter.GET("/admin/timer", GetJudgingTimer)
	adminRouter.GET("/admin/options", GetOptions)
	adminRouter.POST("/admin/options", SetOptions)
	adminRouter.POST("/admin/num-groups", SetNumGroups)
	adminRouter.POST("/admin/group-sizes", SetGroupSizes)
	adminRouter.POST("/admin/block-reqs", SetBlockReqs)
	adminRouter.POST("/admin/max-reqs", SetMaxReqs)

	// Admin panel - exports
	adminRouter.GET("/admin/export/judges", ExportJudges)
	adminRouter.GET("/admin/export/projects", ExportProjects)
	adminRouter.GET("/admin/export/challenges", ExportProjectsByChallenge)
	adminRouter.GET("/admin/export/rankings", ExportRankings)

	// Admin panel - table actions
	adminRouter.PUT("/judge/hide/:id", HideJudge)
	adminRouter.PUT("/project/hide/:id", HideProject)
	adminRouter.PUT("/project/prioritize/:id", PrioritizeProject)
	adminRouter.POST("/project/prioritize", PrioritizeSelectedProjects)
	adminRouter.POST("/project/hide", HideSelectedProjects)
	adminRouter.POST("/judge/hide", HideSelectedJudges)
	adminRouter.POST("/judge/move", MoveSelectedJudges)
	adminRouter.DELETE("/admin/flag/:id", RemoveFlag)
	adminRouter.PUT("/judge/move/:id", MoveJudge)
	adminRouter.PUT("/project/move/:id", MoveProject)
	adminRouter.POST("/project/move", MoveSelectedProjects)
	adminRouter.POST("/admin/deliberation", SetDeliberation)

	// Admin panel - log
	adminRouter.GET("/admin/log", GetLog)

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
	judgeRouter.POST("/judge/break", JudgeBreak)
	judgeRouter.PUT("/judge/notes/:id", JudgeUpdateNotes)
	judgeRouter.GET("/project/:id", GetProject)
	judgeRouter.GET("/project/count", GetProjectCount)
	judgeRouter.GET("/judge/project/:id", GetJudgedProject)
	judgeRouter.GET("/judge/deliberation", GetDeliberationStatus)

	// Project expo routes
	defaultRouter.GET("/project/list/public", ListPublicProjects)
	defaultRouter.GET("/challenges", GetChallenges)
	defaultRouter.GET("/group-names", GetGroupNames)

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

	// Log the successful creation
	logger.SystemLogf("API created and running")

	return router
}

// useVar is a middleware that adds a variable to the context
func useVar(key string, v any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(key, v)
		ctx.Next()
	}
}

// rateLimit is a middleware that limits the number of requests per minute
// for the judge login endpoint
func rateLimit(limiter *Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check for /judge/login endpoint
		if ctx.Request.URL.Path != "/api/judge/login" {
			ctx.Next()
			return
		}

		ip := ctx.ClientIP()
		if !limiter.CheckNewRequest(ip) {
			ctx.AbortWithStatusJSON(429, gin.H{"error": "Too many requests. Logins have been blocked or rate limited."})
			return
		}
		ctx.Next()
	}
}

// getClockFromDb gets the clock state from the database.
// If the clock sync option is not enabled, the clock will be 0.
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

// getLimiterFromDb gets the limiter state from the database
func getLimiterFromDb(db *mongo.Database) *Limiter {
	// Get the limiter state from the database
	options, err := database.GetOptions(db, context.Background())
	if err != nil {
		log.Fatalln("error getting options: " + err.Error())
	}

	limiter := CreateLimiter(int(options.MaxReqPerMin), options.BlockReqs)

	return limiter
}

// Heartbeat is a simple endpoint to check if the server is running
func Heartbeat(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"ok": 1})
}
