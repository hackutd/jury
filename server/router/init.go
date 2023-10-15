package router

import (
	"log"
	"server/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Creates a new Gin router with all routes defined
func NewRouter(db *mongo.Database) *gin.Engine {
	router := gin.Default()

	// Get the clock state from the database
	options, err := database.GetOptions(db)
	if err != nil {
		log.Fatalln("error getting options: " + err.Error())
	}
	clock := options.Clock

	// Add shared variables to router
	router.Use(useVar("db", db))
	router.Use(useVar("clock", &clock))

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

	// Add routes
	judgeRouter.GET("/judge", GetJudge)
	adminRouter.POST("/judge/new", AddJudge)
	defaultRouter.POST("/judge/login", LoginJudge)
	judgeRouter.POST("/judge/auth", JudgeAuthenticated)
	adminRouter.POST("/judge/csv", AddJudgesCsv)
	judgeRouter.GET("/judge/welcome", CheckJudgeReadWelcome)
	judgeRouter.POST("/judge/welcome", SetJudgeReadWelcome)
	adminRouter.GET("/judge/list", ListJudges)
	adminRouter.GET("/judge/stats", JudgeStats)
	adminRouter.DELETE("/judge/:id", DeleteJudge)
	judgeRouter.GET("/judge/projects", GetJudgeProjects)
	judgeRouter.POST("/judge/ipo", GetJudgeIPO)
	judgeRouter.POST("/judge/vote", JudgeVote)
	judgeRouter.POST("/judge/skip", JudgeSkip)
	adminRouter.POST("/project/devpost", AddDevpostCsv)
	adminRouter.POST("/project/new", AddProject)
	adminRouter.GET("/project/list", ListProjects)
	adminRouter.POST("/project/csv", AddProjectsCsv)
	judgeRouter.GET("/project/:id", GetProject)
	judgeRouter.GET("/judge/vote/info", GetVotingProjectInfo)
	judgeRouter.POST("/judge/stars", UpdateStars)
	judgeRouter.GET("/project/count", GetProjectCount)
	judgeRouter.GET("/judge/project/:id", GetJudgedProject)
	adminRouter.DELETE("/project/:id", DeleteProject)
	adminRouter.GET("/project/stats", ProjectStats)
	defaultRouter.POST("/admin/login", LoginAdmin)
	adminRouter.GET("/admin/stats", GetAdminStats)
	adminRouter.GET("/admin/clock", GetClock)
	adminRouter.POST("/admin/clock/pause", PauseClock)
	adminRouter.POST("/admin/clock/unpause", UnpauseClock)
	adminRouter.POST("/admin/clock/reset", ResetClock)
	adminRouter.POST("/admin/auth", AdminAuthenticated)
	adminRouter.POST("/admin/reset", ResetDatabase)
	adminRouter.POST("/judge/hide", HideJudge)
	adminRouter.POST("/judge/unhide", UnhideJudge)
	adminRouter.POST("/project/hide", HideProject)
	adminRouter.POST("/project/unhide", UnhideProject)
	adminRouter.POST("/project/prioritize", PrioritizeProject)
	adminRouter.POST("/project/unprioritize", UnprioritizeProject)
	adminRouter.PUT("/judge/:id", EditJudge)
	defaultRouter.GET("/admin/started", IsClockPaused)
	adminRouter.GET("/admin/flags", GetFlags)

	// Serve frontend static files
	router.Use(static.Serve("/static", static.LocalFile("./public/static", true)))
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.LoadHTMLFiles("./public/index.html")

	// Add no route handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})

	return router
}

func useVar(key string, v any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(key, v)
		ctx.Next()
	}
}
